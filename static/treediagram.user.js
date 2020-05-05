// ==UserScript==
// @name         TreeDiagram
// @namespace    https://www.sshz.org/
// @version      0.1.12.3
// @description  Make Eventernote Great Again
// @author       SSHZ.ORG
// @match        https://www.eventernote.com/*
// @require      https://cdn.jsdelivr.net/npm/chart.js@2.9.3/dist/Chart.bundle.min.js
// @require      https://cdnjs.cloudflare.com/ajax/libs/chartjs-plugin-annotation/0.5.7/chartjs-plugin-annotation.min.js
// @require      https://cdn.jsdelivr.net/npm/chartjs-plugin-zoom@0.7.5/dist/chartjs-plugin-zoom.min.js
// @require      https://cdn.jsdelivr.net/npm/chartjs-plugin-datalabels@0.7.0/dist/chartjs-plugin-datalabels.min.js
// @require      https://cdn.jsdelivr.net/npm/chartjs-plugin-colorschemes@0.4.0/dist/chartjs-plugin-colorschemes.min.js
// ==/UserScript==

(function () {
    'use strict';

    const header = '<h2><ruby>樹形図の設計者<rt>ツリーダイアグラム</rt></ruby></h2>';
    const serverBaseAddress = 'https://treediagram.sshz.org';

    const chartNowAnnotation = {
        type: 'line',
        mode: 'vertical',
        scaleID: 'x-axis-0',
        value: new Date(),
        borderColor: 'rgba(0, 0, 0, 0.2)',
        borderWidth: 1,
        borderDash: [2, 2],
    };

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

        fetch(`${serverBaseAddress}/api/renderEvent?id=${eventId}`)
            .then(response => response.json())
            .then(data => {
                data.snapshots = data.snapshots || [];
                data.snapshots.forEach(snapshot => {
                    snapshot.addedActors = snapshot.addedActors || [];
                    snapshot.removedActors = snapshot.removedActors || [];
                });

                const totalStatsSpan = document.getElementById('td_place_stats_total');
                totalStatsSpan.innerHTML = `${data.placeStatsTotal.rank}/${data.placeStatsTotal.total}`;
                const finishedStatsSpan = document.getElementById('td_place_stats_finished');
                finishedStatsSpan.innerHTML = `${data.placeStatsFinished.rank}/${data.placeStatsFinished.total}`;

                const ctx = document.getElementById('td_chart');
                const annotations = [chartNowAnnotation];

                const liveDate = new Date(data.date);
                liveDate.setUTCHours(3); // JST noon.
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
                            content: 'Live!',
                        },
                    });
                }

                data.snapshots.forEach(snapshot => {
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
                            data: data.snapshots.map(i => {
                                return {
                                    x: new Date(i.timestamp),
                                    y: i.noteCount,
                                };
                            }),
                            cubicInterpolationMode: 'monotone',
                            borderWidth: 1,
                            datalabels: {
                                display: context => data.snapshots[context.dataIndex].dataLabel !== undefined,
                                formatter: (value, context) => data.snapshots[context.dataIndex].dataLabel,
                                align: 'top',
                                backgroundColor: 'rgba(97, 191, 153, 0.5)', // #61BF99
                                borderRadius: 20,
                                color: 'white',
                            },
                        }],
                    },
                    options: {
                        scales: {
                            xAxes: [{
                                type: 'time',
                                ticks: {
                                    maxRotation: 0,
                                },
                                gridLines: {
                                    zeroLineColor: 'rgba(0, 0, 0, 0.1)',
                                },
                            }],
                        },
                        legend: {
                            display: false,
                        },
                        tooltips: {
                            callbacks: {
                                afterLabel: (tooltipItem, chartData) => {
                                    const snapshot = data.snapshots[tooltipItem.index];

                                    let labels = [];
                                    if (snapshot.addedActors.length > 0) {
                                        labels.push('++: ' + snapshot.addedActors.join(', '));
                                    }
                                    if (snapshot.removedActors.length > 0) {
                                        labels.push('--: ' + snapshot.removedActors.join(', '));
                                    }

                                    return labels.join('\n');
                                },
                            },
                        },
                        annotation: { // As of chartjs-plugin-annotation 0.5.7, it does not support `plugins` property.
                            annotations: annotations,
                        },
                        plugins: {
                            zoom: {
                                zoom: {
                                    enabled: true,
                                    drag: true,
                                    mode: 'x',
                                },
                            },
                            colorschemes: {
                                scheme: `tableau.ClassicMedium10`,
                            },
                        },
                    },
                });

                document.getElementById('td_chart_reset').addEventListener('click', () => {
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
            fetch(`${serverBaseAddress}/api/queryEvents?offset=${loadedCount}&${argString}`)
                .then(response => response.json())
                .then(data => {
                    data.forEach(e => {
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

        fetch(`${serverBaseAddress}/api/renderPlace?id=${placeId}`)
            .then(response => response.json())
            .then(data => {
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

        const graphDom = htmlToElement(`
            <div>
                <canvas id="td_chart"></canvas>
                <button class="btn btn-block" type="button" id="td_chart_reset">
                    Reset Zoom
                </button>
            </div>`);
        const favoriteUsersDom = document.getElementsByClassName('gb_users_icon')[0] || document.getElementsByClassName('gb_listusericon')[0];
        favoriteUsersDom.parentNode.insertBefore(graphDom, favoriteUsersDom);

        fetch(`${serverBaseAddress}/api/renderActor?id=${actorId}`)
            .then(response => response.json())
            .then(data => {
                data.snapshots = data.snapshots || [];

                createEventList(`actor=${actorId}`, data.knownEventCount, tdDom, false);

                const ctx = document.getElementById('td_chart');

                data.snapshots.forEach(s => {
                    s.date = new Date(s.date);
                    s.date.setUTCHours(-3); // JST 6am.
                });

                const dataPoints = [];
                if (data.snapshots.length > 0) {
                    const snapshots = data.snapshots.slice();
                    let lastFavoriteCount = 0;

                    for (let date = snapshots[0].date; date < new Date(); date.setDate(date.getDate() + 1)) {
                        if (snapshots.length > 0 && date >= snapshots[0].date) {
                            lastFavoriteCount = snapshots[0].favoriteCount;
                            snapshots.shift();
                        }
                        dataPoints.push({
                            x: new Date(date.getTime()),
                            y: lastFavoriteCount,
                        });
                    }
                }

                const tdChart = new Chart(ctx, {
                    type: 'line',
                    data: {
                        datasets: [{
                            label: 'FavoriteCount',
                            data: dataPoints,
                            cubicInterpolationMode: 'monotone',
                            borderWidth: 1,
                            datalabels: { display: false },
                        }],
                    },
                    options: {
                        scales: {
                            xAxes: [{
                                type: 'time',
                                ticks: {
                                    maxRotation: 0,
                                },
                                gridLines: {
                                    zeroLineColor: 'rgba(0, 0, 0, 0.1)',
                                },
                            }],
                        },
                        legend: {
                            display: false,
                        },
                        annotation: { // As of chartjs-plugin-annotation 0.5.7, it does not support `plugins` property.
                            annotations: [chartNowAnnotation],
                        },
                        plugins: {
                            zoom: {
                                zoom: {
                                    enabled: true,
                                    drag: true,
                                    mode: 'x',
                                },
                            },
                            colorschemes: {
                                scheme: `tableau.ClassicMedium10`,
                            },
                        },
                    },
                });

                document.getElementById('td_chart_reset').addEventListener('click', () => {
                    tdChart.resetZoom();
                });
            });
    }

    function userPage() {
        const favoriteActorsDom = document.getElementsByClassName('gb_actors_list')[0] || document.getElementsByClassName('favorite_actor')[0];
        if (favoriteActorsDom) {
            const actorDoms = favoriteActorsDom.getElementsByTagName('li');
            for (let i = 0; i < actorDoms.length; i++) {
                let count = actorDoms[i].className.match(/c(\d+)/)[1];
                actorDoms[i].getElementsByTagName('a')[0].textContent += ` (${count})`;
            }
        }
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
                                <button class="btn btn-block btn-info" type="button" id="td_compare_actors">Compare Favorites</button>
                                <button class="btn btn-block btn-primary" type="button" id="td_run_query">Query Events</button>
                            </div>
                            <div id="td_actor_favorites_chart_container"></div>
                            <div id="td_event_list_container"></div>
                        </div>
                    </div>
                </div>
            </div>`);
        const contentDom = document.getElementsByClassName('gb_ad_footer').length > 0 ? document.getElementsByClassName('gb_ad_footer')[0].previousElementSibling : document.getElementById('container');
        contentDom.replaceWith(tdDom);

        const selectedActors = [];
        const actorNames = {};
        const selectedActorsDom = document.getElementById('td_selected_actors');
        const searchActorResultDom = document.getElementById('td_search_actor_result');
        const actorFavoritesChartContainerDom = document.getElementById('td_actor_favorites_chart_container');
        const eventListContainerDom = document.getElementById('td_event_list_container');

        const inputDom = document.getElementById('td_search_actor_input');
        inputDom.addEventListener('keyup', () => {
            searchActor(inputDom.value);
        });

        document.getElementById('td_compare_actors').addEventListener('click', compareActors);
        document.getElementById('td_run_query').addEventListener('click', runQuery);

        function addActor(id, name) {
            if (selectedActors.some(i => i === id)) return;
            selectedActors.push(id);
            actorNames[id] = name;

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

        function compareActors() {
            while (actorFavoritesChartContainerDom.firstChild) {
                actorFavoritesChartContainerDom.removeChild(actorFavoritesChartContainerDom.firstChild);
            }

            const ctx = htmlToElement(`<canvas id="td_chart"></canvas>`);
            actorFavoritesChartContainerDom.appendChild(ctx);

            fetch(`${serverBaseAddress}/api/compareActors?${selectedActors.map(i => `id=${i}`).join('&')}`)
                .then(response => response.json())
                .then(data => {
                    const datasets = [];

                    for (const [actorId, res] of Object.entries(data)) {
                        res.snapshots.forEach(s => {
                            s.date = new Date(s.date);
                            s.date.setUTCHours(-3); // JST 6am.
                        });

                        const dataPoints = [];
                        if (res.snapshots.length > 0) {
                            const snapshots = res.snapshots.slice();
                            let lastFavoriteCount = 0;

                            for (let date = snapshots[0].date; date < new Date(); date.setDate(date.getDate() + 1)) {
                                if (snapshots.length > 0 && date >= snapshots[0].date) {
                                    lastFavoriteCount = snapshots[0].favoriteCount;
                                    snapshots.shift();
                                }
                                dataPoints.push({
                                    x: new Date(date.getTime()),
                                    y: lastFavoriteCount,
                                });
                            }
                        }

                        datasets.push(
                            {
                                label: `${actorNames[actorId]}`,
                                data: dataPoints,
                                cubicInterpolationMode: 'monotone',
                                borderWidth: 1,
                                datalabels: { display: false },
                            }
                        );
                    }

                    datasets.sort(function (a, b) {
                        return a.data[a.data.length - 1].y - b.data[b.data.length - 1].y;
                    });

                    const tdChart = new Chart(ctx, {
                        type: 'line',
                        data: {
                            datasets: datasets,
                        },
                        options: {
                            scales: {
                                xAxes: [{
                                    type: 'time',
                                    ticks: {
                                        maxRotation: 0,
                                    },
                                    gridLines: {
                                        zeroLineColor: 'rgba(0, 0, 0, 0.1)',
                                    },
                                }],
                            },
                            tooltips: {
                                mode: 'index',
                            },
                            annotation: { // As of chartjs-plugin-annotation 0.5.7, it does not support `plugins` property.
                                annotations: [chartNowAnnotation],
                            },
                            plugins: {
                                colorschemes: {
                                    scheme: `brewer.Spectral${Math.min(11, Math.max(3, datasets.length))}`,
                                },
                            },
                        },
                    });
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

        const eventPageRegex = /^https:\/\/www.eventernote.com\/events\/(\d+)$/;
        const eventPageMatch = url.match(eventPageRegex);
        if (eventPageMatch) {
            return eventPage(eventPageMatch[1]);
        }

        const placePageRegex = /^https:\/\/www.eventernote.com\/places\/(\d+)$/;
        const placePageMatch = url.match(placePageRegex);
        if (placePageMatch) {
            return placePage(placePageMatch[1]);
        }

        const actorPageRegex = /^https:\/\/www.eventernote.com\/actors\/(?:.+\/)?(\d+)$/;
        const actorPageMatch = url.match(actorPageRegex);
        if (actorPageMatch) {
            return actorPage(actorPageMatch[1]);
        }

        const userPageRegex = /^https:\/\/www.eventernote.com\/users(?:\/.+)?$/;
        const userPageMatch = url.match(userPageRegex);
        if (userPageMatch) {
            return userPage();
        }
    }

    main();
})();
