// ==UserScript==
// @name         TreeDiagram
// @namespace    https://www.sshz.org/
// @version      0.1.10.2
// @description  Make Eventernote Great Again
// @author       SSHZ.ORG
// @match        https://www.eventernote.com/*
// @require      https://cdn.jsdelivr.net/npm/chart.js@2.8.0/dist/Chart.bundle.min.js
// @require      https://cdnjs.cloudflare.com/ajax/libs/chartjs-plugin-annotation/0.5.7/chartjs-plugin-annotation.min.js
// @require      https://cdn.jsdelivr.net/npm/chartjs-plugin-zoom@0.7.0/dist/chartjs-plugin-zoom.min.js
// @require      https://cdn.jsdelivr.net/npm/chartjs-plugin-datalabels@0.6.0/dist/chartjs-plugin-datalabels.min.js
// ==/UserScript==

(function () {
    'use strict';

    const header = '<h2><ruby>樹形図の設計者<rt>ツリーダイアグラム</rt></ruby></h2>';
    const serverBaseAddress = 'https://treediagram.sshz.org';

    function htmlToElement(html) {
        const template = document.createElement('template');
        html = html.trim();
        template.innerHTML = html;
        return template.content.firstChild;
    }

    function eventPage(eventId) {
        const tdDom = htmlToElement(`
            <div>
                ${header}
                <canvas id="td_chart"></canvas>
                <button class="btn btn-block" type="button" id="td_chart_reset">
                    Reset Zoom
                </button>
                <div class="ma10 alert alert-info s">
                    Rank in Place:
                    <span class="label">Total</span>
                    <span id="td_place_stats_total"></span>
                    <span class="label">Finished</span>
                    <span id="td_place_stats_finished"></span>
                </div>
            </div>`);

        const entryAreaDom = document.getElementById('entry_area') || document.getElementsByClassName('mod_events_detail')[0];
        entryAreaDom.parentNode.insertBefore(tdDom, entryAreaDom.nextSibling);

        fetch(`${serverBaseAddress}/api/renderEvent?id=${eventId}`).then(function (response) {
            return response.json();
        }).then(function (data) {
            const totalStatsSpan = document.getElementById("td_place_stats_total");
            totalStatsSpan.innerHTML = `${data.placeStatsTotal.rank}/${data.placeStatsTotal.total}`;
            const finishedStatsSpan = document.getElementById("td_place_stats_finished");
            finishedStatsSpan.innerHTML = `${data.placeStatsFinished.rank}/${data.placeStatsFinished.total}`;

            const ctx = document.getElementById("td_chart");
            const annotations = [{
                type: 'line',
                mode: 'vertical',
                scaleID: 'x-axis-0',
                value: new Date(),
                borderColor: 'rgba(0, 0, 0, 0.2)',
                borderWidth: 1,
                borderDash: [2, 2],
            }];

            const liveDate = new Date(data.date);
            liveDate.setUTCHours(3);  // JST noon.
            if (data.snapshots.length > 0 && new Date(data.snapshots[0].timestamp) <= liveDate) {
                annotations.push({
                    type: 'line',
                    mode: 'vertical',
                    scaleID: 'x-axis-0',
                    value: liveDate,
                    borderColor: 'rgba(0, 0, 0, 0.2)',
                    borderWidth: 1,
                    borderDash: [10, 10],
                    label: {
                        backgroundColor: 'rgba(0, 0, 0, 0.4)',
                        cornerRadius: 0,
                        position: 'bottom',
                        enabled: true,
                        content: "Live!"
                    }
                });
            }

            data.snapshots.forEach(function (snapshot) {
                if (snapshot.addedActors.length > 0 || snapshot.removedActors.length > 0) {
                    let label = '';
                    if (snapshot.addedActors.length > 0) {
                        label += '+';
                    }
                    if (snapshot.removedActors.length > 0) {
                        label += '-';
                    }
                    snapshot.dataLabel = label;
                }
            });

            const tdChart = new Chart(ctx, {
                type: 'line',
                data: {
                    datasets: [{
                        label: 'NoteCount',
                        data: data.snapshots.map(function (i) {
                            return {
                                x: new Date(i.timestamp),
                                y: i.noteCount,
                                snapshot: i,
                            };
                        }),
                        cubicInterpolationMode: 'monotone',
                        backgroundColor: 'rgba(54, 162, 235, 0.5)',
                        borderColor: 'rgb(54, 162, 235)',
                        borderWidth: 1,
                        datalabels: {
                            display: function (context) {
                                return data.snapshots[context.dataIndex].dataLabel !== undefined;
                            },
                            formatter: function (value, context) {
                                return data.snapshots[context.dataIndex].dataLabel;
                            },
                            align: 'top',
                            backgroundColor: 'rgba(97, 191, 153, 0.5)',
                            borderRadius: 20,
                            color: 'white'
                        }
                    }]
                },
                options: {
                    scales: {
                        xAxes: [{
                            type: 'time',
                            ticks: {
                                maxRotation: 0
                            }
                        }]
                    },
                    legend: {
                        display: false
                    },
                    tooltips: {
                        callbacks: {
                            afterLabel: function (tooltipItem, data) {
                                const snapshot = data.datasets[tooltipItem.datasetIndex].data[tooltipItem.index].snapshot;

                                let labels = [];
                                if (snapshot.addedActors.length > 0) {
                                    labels.push('++: ' + snapshot.addedActors.join(', '));
                                }
                                if (snapshot.removedActors.length > 0) {
                                    labels.push('--: ' + snapshot.removedActors.join(', '));
                                }

                                return labels.join('\n');
                            }
                        }
                    },
                    annotation: { // As of chartjs-plugin-annotation 0.5.7, it does not support `plugins` property.
                        annotations: annotations
                    },
                    plugins: {
                        zoom: {
                            zoom: {
                                enabled: true,
                                drag: true,
                                mode: 'x',
                            },
                        },
                    },
                },
            });

            document.getElementById('td_chart_reset').addEventListener('click', function () {
                tdChart.resetZoom();
            });
        });
    }

    function createEventList(argString, totalCount, domToAppend, autoLoadFirstPage) {
        let loadedCount = 0;

        const eventListDom = htmlToElement(`
            <div>
                <table class="table table-hover s" style="margin-bottom: 0;">
                    <tbody id="td_event_list_tbody"></tbody>
                </table>
                <button class="btn btn-block" type="button" id="td_event_list_load_more_button" disabled>
                    Load More (<span id="td_event_list_loaded_indicator">0</span> / ${totalCount || 'Unknown'})
                </button>
            </div>`);
        domToAppend.appendChild(eventListDom);

        const tbodyDom = document.getElementById('td_event_list_tbody');
        const loadMoreButtonDom = document.getElementById('td_event_list_load_more_button');
        const loadedIndicatorDom = document.getElementById('td_event_list_loaded_indicator');

        function loadMore() {
            if (totalCount && loadedCount >= totalCount) {
                return;
            }

            loadMoreButtonDom.disabled = true;
            fetch(`${serverBaseAddress}/api/queryEvents?offset=${loadedCount}&${argString}`).then(function (response) {
                return response.json();
            }).then(function (data) {
                data.forEach(function (e) {
                    const trDom = htmlToElement(`
                    <tr>
                        <td>${loadedCount + 1}</td>
                        <td nowrap>${e.date}</td>
                        <td><a href="/events/${e.id}" target="_blank">${e.name}</a></td>
                        <td>${e.lastNoteCount}</td>
                    </tr>`);

                    if (!e.finished) {
                        trDom.classList.add('warning');
                    }

                    tbodyDom.appendChild(trDom);
                    loadedCount += 1;
                });

                loadedIndicatorDom.innerText = loadedCount.toString();

                if ((!totalCount) || loadedCount < totalCount) {
                    loadMoreButtonDom.disabled = false;
                }
            });
        }

        loadMoreButtonDom.addEventListener('click', loadMore);
        if ((!totalCount) || totalCount > 0) {
            loadMoreButtonDom.disabled = false;
        }

        if (autoLoadFirstPage) {
            loadMore();
        }
    }

    function placePage(placeId) {
        const tdDom = htmlToElement(`
            <div>
                ${header}
                <h3>Top Event List</h3>
            </div>`);

        const placeDetailsTableDom = document.getElementsByClassName('gb_place_detail_table')[0] || document.getElementsByClassName('mod_places_detail')[0];
        placeDetailsTableDom.parentNode.insertBefore(tdDom, placeDetailsTableDom.nextSibling);

        fetch(`${serverBaseAddress}/api/renderPlace?id=${placeId}`).then(function (response) {
            return response.json();
        }).then(function (data) {
            createEventList(`place=${placeId}`, data.knownEventCount, tdDom, false);
        });
    }

    function actorPage(actorId) {
        const tdDom = htmlToElement(`
            <div>
                ${header}
                <h3>Top Event List</h3>
            </div>`);

        const actorTitleDom = document.getElementsByClassName('gb_actors_title')[0] || document.getElementsByClassName('gb_blur_title')[0];
        actorTitleDom.parentNode.insertBefore(tdDom, actorTitleDom.nextSibling);

        fetch(`${serverBaseAddress}/api/renderActor?id=${actorId}`).then(function (response) {
            return response.json();
        }).then(function (data) {
            createEventList(`actor=${actorId}`, data.knownEventCount, tdDom, false);
        });
    }

    function treeDiagramPage() {
        const tdDom = htmlToElement(`
            <div class="container">
                <div class="row">
                    <div class="page span8">
                        <div class="page-header">${header}</div>
                        <div>
                            <h2>Advanced Event Search</h2>
                            <div class="gb_form">
                            <table class="table table-bordered table-striped">
                            <tbody>
                            <tr>
                                <td class="span2">Actors</td>
                                <td>
                                    <div id="td_selected_actors"></div>
                                </td>
                            </tr>
                            <tr>
                                <td></td>
                                <td>
                                    <input type="text" id="td_search_actor_input" placeholder="Search actors..." />
                                    <div id="td_search_actor_result"></div>
                                </td>
                            </tr>
                            </tbody>
                            </table>
                            </div>
                            <div class="form-actions">
                                <button class="btn btn-block btn-primary" type="button" id="td_run_query">Query</button>
                            </div>
                            <div id="td_event_list_container"></div>
                        </div>
                    </div>
                </div>
            </div>`);
        const contentDom = document.getElementsByClassName('gb_ad_footer').length > 0 ? document.getElementsByClassName('gb_ad_footer')[0].previousElementSibling : document.getElementById('container');
        contentDom.replaceWith(tdDom);

        const selectedActors = [];
        const selectedActorsDom = document.getElementById('td_selected_actors');
        const searchActorResultDom = document.getElementById('td_search_actor_result');
        const eventListContainerDom = document.getElementById('td_event_list_container');

        const inputDom = document.getElementById('td_search_actor_input');
        inputDom.addEventListener('keyup', () => {
            searchActor(inputDom.value);
        });

        document.getElementById('td_run_query').addEventListener('click', runQuery);

        function addActor(id, name) {
            if (selectedActors.some(i => i === id)) return;
            selectedActors.push(id);

            const buttonDom = htmlToElement(`
                <button class="btn" type="button"><i class="icon-minus"></i> ${name}</button>`);
            selectedActorsDom.appendChild(buttonDom);

            buttonDom.addEventListener('click', () => {
                removeActor(id, buttonDom);
            });
        }

        function removeActor(id, liDom) {
            const index = selectedActors.indexOf(id);
            if (index > -1) {
                selectedActors.splice(index, 1);
                liDom.remove();
            }
        }

        function searchActor(keyword) {
            if (!keyword) return;
            fetch(`/api/actors/search?&simple=1&limit=20&keyword=${keyword}`)
                .then(response => response.json())
                .then(data => {
                    while (searchActorResultDom.firstChild) {
                        searchActorResultDom.removeChild(searchActorResultDom.firstChild);
                    }
                    if (data.results) {
                        data.results.forEach(item => {
                            const itemDom = htmlToElement(`
                                <button class="btn" type="button"><i class="icon-plus"></i> ${item.name}</button>`);
                            itemDom.addEventListener('click', () => {
                                addActor(item.id, item.name);
                            });
                            searchActorResultDom.appendChild(itemDom);
                        });
                    }
                });
        }

        function runQuery() {
            if (eventListContainerDom.firstChild) {
                eventListContainerDom.removeChild(eventListContainerDom.firstChild);
            }
            createEventList(selectedActors.map(i => `actor=${i}`).join('&'), undefined, eventListContainerDom, true);
        }
    }

    function addTreeDiagramPageLink() {
        const liDom = htmlToElement(`<li><a href="#">TreeDiagram</a></li>`);
        liDom.addEventListener('click', () => {
            treeDiagramPage();
        });

        const nav = document.getElementsByClassName('nav')[0] || document.getElementsByClassName('grid')[0].firstElementChild;
        nav.appendChild(liDom);
    }

    function main() {
        addTreeDiagramPageLink();

        const url = document.URL;

        const eventPageRegex = /https:\/\/www.eventernote.com\/events\/(\d+)/g;
        const eventPageMatch = eventPageRegex.exec(url);
        if (eventPageMatch) {
            return eventPage(eventPageMatch[1]);
        }

        const placePageRegex = /https:\/\/www.eventernote.com\/places\/(\d+)/g;
        const placePageMatch = placePageRegex.exec(url);
        if (placePageMatch) {
            return placePage(placePageMatch[1]);
        }

        const actorPageRegex = /https:\/\/www.eventernote.com\/actors\/(?:.+\/)?(\d+)/g;
        const actorPageMatch = actorPageRegex.exec(url);
        if (actorPageMatch) {
            return actorPage(actorPageMatch[1]);
        }
    }

    main();
})();
