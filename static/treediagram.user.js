// ==UserScript==
// @name         TreeDiagram
// @namespace    https://www.sshz.org/
// @version      0.1.9.3
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

        const entryAreaDom = document.getElementById('entry_area');
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
                    var label = '';
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

    function createEventList(argString, totalCount, domToAppend) {
        let loadedCount = 0;

        const eventListDom = htmlToElement(`
            <div>
                <table class="table table-hover s" style="margin-bottom: 0;">
                    <tbody id="td_event_list_tbody"></tbody>
                </table>
                <button class="btn btn-block" type="button" id="td_event_list_load_more_button" disabled>
                    Load More (<span id="td_event_list_loaded_indicator">0</span> / ${totalCount})
                </button>
            </div>`);
        domToAppend.appendChild(eventListDom);

        const tbodyDom = document.getElementById('td_event_list_tbody');
        const loadMoreButtonDom = document.getElementById('td_event_list_load_more_button');
        const loadedIndicatorDom = document.getElementById('td_event_list_loaded_indicator');

        function loadMore() {
            if (loadedCount >= totalCount) {
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
                        <td><a href="/events/${e.id}">${e.name}</a></td>
                        <td>${e.lastNoteCount}</td>
                    </tr>`);

                    if (!e.finished) {
                        trDom.classList.add('warning');
                    }

                    tbodyDom.appendChild(trDom);
                    loadedCount += 1;
                });

                loadedIndicatorDom.innerText = loadedCount.toString();

                if (loadedCount < totalCount) {
                    loadMoreButtonDom.disabled = false;
                }
            });
        }

        loadMoreButtonDom.addEventListener('click', loadMore);
        if (totalCount > 0) {
            loadMoreButtonDom.disabled = false;
        }
    }

    function placePage(placeId) {
        const tdDom = htmlToElement(`
            <div>
                ${header}
                <h3>Top Event List</h3>
            </div>`);

        const placeDetailsTableDom = document.getElementsByClassName('gb_place_detail_table')[0];
        placeDetailsTableDom.parentNode.insertBefore(tdDom, placeDetailsTableDom.nextSibling);

        fetch(`${serverBaseAddress}/api/renderPlace?id=${placeId}`).then(function (response) {
            return response.json();
        }).then(function (data) {
            createEventList(`place=${placeId}`, data.knownEventCount, tdDom);
        });
    }

    function actorPage(actorId) {
        const tdDom = htmlToElement(`
            <div>
                ${header}
                <h3>Top Event List</h3>
            </div>`);

        const actorTitleDom = document.getElementsByClassName('gb_actors_title')[0];
        actorTitleDom.parentNode.insertBefore(tdDom, actorTitleDom.nextSibling);

        fetch(`${serverBaseAddress}/api/renderActor?id=${actorId}`).then(function (response) {
            return response.json();
        }).then(function (data) {
            createEventList(`actor=${actorId}`, data.knownEventCount, tdDom);
        });
    }

    function main() {
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
