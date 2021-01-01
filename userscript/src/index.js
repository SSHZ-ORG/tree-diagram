const servicepb = require('./service_grpc_web_pb');
const pb = require('./service_pb');

let Highcharts = require('highcharts');
require('highcharts/modules/exporting')(Highcharts);

const treeDiagramService = new servicepb.TreeDiagramServicePromiseClient('https://treediagram.sshz.org');

const header = '<h2><ruby>樹形図の設計者<rt>ツリーダイアグラム</rt></ruby></h2>';

const plotLineNow = {
    value: new Date(),
    dashStyle: 'Dash',
    label: {text: 'Now'},
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
            <div id="td_chart"></div>
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

    const request = new pb.RenderEventRequest().setId(eventId);
    treeDiagramService.renderEvent(request).then(response => {
        const totalStatsSpan = document.getElementById('td_place_stats_total');
        totalStatsSpan.innerHTML = `${response.getPlaceStatsTotal().getRank()}/${response.getPlaceStatsTotal().getTotal()}`;
        const finishedStatsSpan = document.getElementById('td_place_stats_finished');
        finishedStatsSpan.innerHTML = `${response.getPlaceStatsFinished().getRank()}/${response.getPlaceStatsFinished().getTotal()}`;

        const ctx = document.getElementById('td_chart');
        const plotLines = [plotLineNow];

        const liveDate = new Date(response.getDate());
        liveDate.setUTCHours(3); // JST noon.
        if (response.getSnapshotsList().length > 0 && response.getSnapshotsList()[0].getTimestamp().toDate() <= liveDate) {
            plotLines.push({
                value: liveDate,
                dashStyle: 'LongDashDot',
                label: {text: 'Live!'},
            });
        }

        Highcharts.chart(ctx, {
            chart: {zoomType: 'x'},
            credits: {enabled: false},
            title: {text: undefined},
            plotOptions: {areaspline: {threshold: null}},
            xAxis: {
                type: 'datetime',
                plotLines: plotLines,
            },
            yAxis: {title: {text: undefined}},
            legend: {enabled: false},
            series: [
                {
                    name: 'NoteCount',
                    type: 'areaspline',
                    data: response.getSnapshotsList().map(i => {
                        const p = {x: i.getTimestamp().toDate(), y: i.getNoteCount()};
                        if (i.getAddedActorsList().length > 0 || i.getRemovedActorsList().length > 0) {
                            let label = '';
                            if (i.getAddedActorsList().length > 0) {
                                label += '+';
                            }
                            if (i.getRemovedActorsList().length > 0) {
                                label += '-';
                            }
                            p.dataLabels = {
                                enabled: true,
                                color: 'white',
                                backgroundColor: 'rgba(97, 191, 153, 0.5)', // #61BF99
                                borderRadius: 5,
                                format: label,
                                crop: false,
                                overflow: 'allow',
                            };
                        }
                        return p;
                    }),
                    tooltip: {
                        pointFormatter: function () {
                            let labels = [`NoteCount: <b>${this.y}</b>`];
                            const snapshot = response.getSnapshotsList()[this.index];

                            if (snapshot.getAddedActorsList().length > 0) {
                                labels.push('++: ' + snapshot.getAddedActorsList().join(', '));
                            }
                            if (snapshot.getRemovedActorsList().length > 0) {
                                labels.push('--: ' + snapshot.getRemovedActorsList().join(', '));
                            }
                            return labels.join('<br/>');
                        },
                    },
                },
                {
                    // Dummy series to make sure plotLines appear even if they are out of main data range.
                    type: 'scatter',
                    marker: {enabled: false},
                    data: plotLines.map(i => ({x: i.value})),
                },
            ],
        });
    });
}

function createEventList(placeId, actorIds, totalCount, domToAppend, autoLoadFirstPage) {
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
        const request = new pb.QueryEventsRequest().setOffset(loadedCount);
        if (placeId !== "") {
            request.setPlaceId(placeId);
        }
        actorIds.forEach(i => request.addActorIds(i.toString()));
        treeDiagramService.queryEvents(request).then(response => {
            response.getEventsList().forEach(e => {
                const trDom = htmlToElement(`
                    <tr>
                        <td>${loadedCount + 1}</td>
                        <td nowrap>${e.getDate()}</td>
                        <td><a href="/events/${e.getId()}" target="_blank">${e.getName()}</a></td>
                        <td>${e.getLastNoteCount()}</td>
                    </tr>`);

                if (!e.getFinished()) {
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

    const request = new pb.RenderPlaceRequest().setId(placeId);
    treeDiagramService.renderPlace(request).then(response => {
        createEventList(placeId, [], response.getKnownEventCount(), tdDom, false);
    });
}

function actorSnapshotsToDataPoints(snapshots) {
    if (snapshots.length === 0) {
        return [];
    }

    const dates = snapshots.map(s => {
        const d = new Date(s.getDate());
        d.setUTCHours(-3); // JST 6am.
        return d;
    });

    const dataPoints = [];
    const copy = snapshots.slice();
    let lastFavoriteCount = 0;

    for (let date = dates[0]; date < new Date(); date.setDate(date.getDate() + 1)) {
        if (copy.length > 0 && date >= dates[0]) {
            lastFavoriteCount = copy[0].getFavoriteCount();
            dates.shift();
            copy.shift();
        }
        dataPoints.push({
            x: date.getTime(),
            y: lastFavoriteCount,
        });
    }
    return dataPoints;
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
            <div id="td_chart"></div>
        </div>`);
    const favoriteUsersDom = document.getElementsByClassName('gb_users_icon')[0] || document.getElementsByClassName('gb_listusericon')[0];
    favoriteUsersDom.parentNode.insertBefore(graphDom, favoriteUsersDom);

    const request = new pb.RenderActorsRequest().addId(actorId);
    treeDiagramService.renderActors(request).then(response => {
        const data = response.getItemsMap().get(actorId);

        createEventList("", [actorId], data.getKnownEventCount(), tdDom, false);

        const ctx = document.getElementById('td_chart');

        const plotLines = [plotLineNow];
        Highcharts.chart(ctx, {
            chart: {zoomType: 'x'},
            credits: {enabled: false},
            title: {text: undefined},
            plotOptions: {areaspline: {threshold: null}},
            xAxis: {
                type: 'datetime',
                plotLines: plotLines,
            },
            yAxis: {title: {text: undefined}},
            legend: {enabled: false},
            tooltip: {xDateFormat: '%Y-%m-%d'},
            series: [
                {
                    name: 'FavoriteCount',
                    type: 'areaspline',
                    data: actorSnapshotsToDataPoints(data.getSnapshotsList()),
                },
                {
                    // Dummy series to make sure plotLines appear even if they are out of main data range.
                    type: 'scatter',
                    marker: {enabled: false},
                    data: plotLines.map(i => ({x: i.value})),
                },
            ],
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

        actorFavoritesChartContainerDom.appendChild(htmlToElement(`<div id="td_chart"></div>`));
        const ctx = document.getElementById('td_chart');

        const request = new pb.RenderActorsRequest();
        selectedActors.forEach(i => request.addId(i.toString()));

        treeDiagramService.renderActors(request).then(response => {
            const series = [];
            for (const [actorId, res] of response.getItemsMap().entries()) {
                series.push({
                    name: `${actorNames[actorId]}`,
                    type: 'spline',
                    data: actorSnapshotsToDataPoints(res.getSnapshotsList()),
                });
            }
            series.sort(function (a, b) {
                return a.data[a.data.length - 1].y - b.data[b.data.length - 1].y;
            });
            series.reverse();

            Highcharts.chart(ctx, {
                chart: {zoomType: 'x'},
                credits: {enabled: false},
                title: {text: undefined},
                plotOptions: {spline: {threshold: null}},
                xAxis: {type: 'datetime'},
                yAxis: {title: {text: undefined}},
                series: series,
                tooltip: {xDateFormat: '%Y-%m-%d', shared: true},
            });
        });
    }

    function runQuery() {
        if (eventListContainerDom.firstChild) {
            eventListContainerDom.removeChild(eventListContainerDom.firstChild);
        }
        createEventList("", selectedActors, undefined, eventListContainerDom, true);
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
