import * as pb from "./service_pb";
import Highcharts, {Options, PointOptionsObject} from "highcharts";
import {TreeDiagramServiceClient} from "./ServiceServiceClientPb";

(function () {
    'use strict';

    Highcharts.setOptions({
        chart: {
            panning: {enabled: true},
            panKey: "alt",
        },
        credits: {enabled: false},
        title: {text: undefined},
        plotOptions: {
            areaspline: {threshold: null},
            spline: {threshold: null},
        },
        time: {timezoneOffset: new Date().getTimezoneOffset()},
        accessibility: {enabled: false},
    });

    const treeDiagramService = new TreeDiagramServiceClient('https://treediagram.sshz.org');
    const renderActorsMax = 100;

    const header = '<h2><ruby>樹形図の設計者<rt>ツリーダイアグラム</rt></ruby></h2>';

    function formatDate(date: pb.Date): string {
        return `${date.getYear()}-${date.getMonth().toString().padStart(2, '0')}-${date.getDay().toString().padStart(2, '0')}`;
    }

    function convertToJsDate(date: pb.Date, utcHours: number): Date {
        // This explodes if date.year is between 0 and 99.
        // Empty value in BE can cause year to be 0001, so we will wrongly think it's some time in 1901 here.
        return new Date(Date.UTC(date.getYear(), date.getMonth() - 1, date.getDay(), utcHours));
    }

    function htmlToElement(html: string) {
        const template = document.createElement('template');
        html = html.trim();
        template.innerHTML = html;
        return template.content.firstChild as HTMLElement;
    }

    function eventPage(eventId: string) {
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

        const actorsUlDom = document.getElementsByClassName("actors");
        const actorNames = new Map<string, string>();
        const isDesktopLayout = actorsUlDom.length > 0;

        if (isDesktopLayout) {
            for (let liDom of actorsUlDom[0].children) {
                const aDom = liDom.children[0] as HTMLAnchorElement;
                const actorId = aDom.href.substring(aDom.href.lastIndexOf("/") + 1);
                actorNames.set(actorId, aDom.text);
            }
        }

        const request = new pb.RenderEventRequest().setId(eventId);
        treeDiagramService.renderEvent(request, null).then(response => {
            const totalStatsSpan = document.getElementById('td_place_stats_total');
            totalStatsSpan.innerHTML = `${response.getPlaceStatsTotal().getRank()}/${response.getPlaceStatsTotal().getTotal()}`;
            const finishedStatsSpan = document.getElementById('td_place_stats_finished');
            finishedStatsSpan.innerHTML = `${response.getPlaceStatsFinished().getRank()}/${response.getPlaceStatsFinished().getTotal()}`;

            const plotLines = [{
                value: new Date().getTime(),
                dashStyle: 'Dash',
                label: {text: 'Now'},
            }];

            const compressedSnapshots = response.getCompressedSnapshotsList();

            const liveDate = convertToJsDate(response.getDate(), 3); // JST noon.
            if (compressedSnapshots.length > 0 && compressedSnapshots[0].getTimestampsList()[0].toDate() <= liveDate) {
                plotLines.push({
                    value: liveDate.getTime(),
                    dashStyle: 'LongDashDot',
                    label: {text: 'Live!'},
                });
            }

            const snapshots: { time: number, noteCount: number, addedActors: string[], removedActors: string[] }[] = [];
            for (let snapshot of compressedSnapshots) {
                for (let i = 0; i < snapshot.getTimestampsList().length; i++) {
                    const timestamp = snapshot.getTimestampsList()[i];
                    snapshots.push({
                        time: timestamp.toDate().getTime(),
                        noteCount: snapshot.getNoteCount(),
                        addedActors: i === 0 ? snapshot.getAddedActorsList().map(a => a.getName()) : [],
                        removedActors: i === 0 ? snapshot.getRemovedActorsList().map(a => a.getName()) : [],
                    });
                }

                if (isDesktopLayout) {
                    [...snapshot.getAddedActorsList(), ...snapshot.getRemovedActorsList()].forEach(a => {
                        if (actorNames.has(a.getId())) return;
                        actorNames.set(a.getId(), a.getName());
                        actorsUlDom[0].appendChild(htmlToElement(`
                            <li><a href="/actors/${a.getName()}/${a.getId()}">${a.getName()}</a>*</li>`));
                    });
                }
            }

            Highcharts.chart('td_chart', {
                chart: {zoomType: 'x'},
                xAxis: {
                    type: 'datetime',
                    plotLines: plotLines,
                },
                yAxis: {title: {text: undefined}},
                legend: {enabled: false},
                tooltip: {xDateFormat: '%Y-%m-%d %H:%M:%S.%L'},
                series: [
                    {
                        name: 'NoteCount',
                        type: 'areaspline',
                        data: snapshots.map(snapshot => {
                            const p: PointOptionsObject = {x: snapshot.time, y: snapshot.noteCount};
                            if (snapshot.addedActors.length > 0 || snapshot.removedActors.length > 0) {
                                let label = '';
                                if (snapshot.addedActors.length > 0) {
                                    label += '+';
                                }
                                if (snapshot.removedActors.length > 0) {
                                    label += '-';
                                }
                                p.dataLabels = {
                                    enabled: true,
                                    color: 'white',
                                    backgroundColor: 'rgba(97, 191, 153, 0.5)', // #61BF99
                                    borderRadius: 5,
                                    format: label,
                                    allowOverlap: true,
                                    crop: false,
                                    overflow: 'allow',
                                };
                            }
                            return p;
                        }),
                        tooltip: {
                            pointFormatter: function () {
                                let labels = [`NoteCount: <b>${this.y}</b>`];
                                const snapshot = snapshots[this.index];

                                if (snapshot.addedActors.length > 0) {
                                    labels.push('++: ' + snapshot.addedActors.join(', '));
                                }
                                if (snapshot.removedActors.length > 0) {
                                    labels.push('--: ' + snapshot.removedActors.join(', '));
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
            } as Options);

            if (isDesktopLayout && actorNames.size <= renderActorsMax) {
                actorsUlDom[0].parentElement.appendChild(htmlToElement(`
                    <div>
                        <button class="btn btn-block" type="button" id="td_event_compare_actors">Compare Favorites (${actorNames.size})</button>
                        <div id="td_event_actor_favorites_chart_container"></div>
                    </div>`));

                const buttonEl = document.getElementById("td_event_compare_actors") as HTMLButtonElement;
                buttonEl.addEventListener('click', () => {
                    buttonEl.disabled = true;
                    compareActors(document.getElementById("td_event_actor_favorites_chart_container"),
                        actorNames, () => buttonEl.remove(),
                        [{label: "This Event", time: convertToJsDate(response.getDate(), 3)}]);
                });
            }
        });
    }

    function createEventList(filter: pb.QueryEventsRequest.EventFilter, totalCount: number | undefined, domToAppend: HTMLElement, autoLoadFirstPage: boolean) {
        let loadedCount = 0;

        const eventListDom = htmlToElement(`
        <div>
            <table class="table table-hover s" style="margin-bottom: 0;">
                <tbody id="td_event_list_tbody"></tbody>
            </table>
            <button class="btn btn-block" type="button" id="td_event_list_load_more_button" disabled>
                Load More (<span id="td_event_list_loaded_indicator">0</span> / <span id="td_event_list_total_indicator">${totalCount === undefined ? 'Unknown' : totalCount}</span>)
            </button>
        </div>`);
        domToAppend.appendChild(eventListDom);

        const tbodyDom = document.getElementById('td_event_list_tbody');
        const loadMoreButtonDom = document.getElementById('td_event_list_load_more_button') as HTMLButtonElement;
        const loadedIndicatorDom = document.getElementById('td_event_list_loaded_indicator');
        const totalIndicatorDom = document.getElementById('td_event_list_total_indicator');

        function enableLoadMoreButton() {
            if (totalCount === undefined || totalCount > loadedCount) {
                loadMoreButtonDom.disabled = false;
            }
        }

        function loadMore() {
            if (totalCount !== undefined && loadedCount >= totalCount) {
                return;
            }

            loadMoreButtonDom.disabled = true;
            const request = new pb.QueryEventsRequest().setOffset(loadedCount).setFilter(filter);
            treeDiagramService.queryEvents(request, null).then(response => {
                response.getEventsList().forEach(e => {
                    const trDom = htmlToElement(`
                    <tr>
                        <td>${loadedCount + 1}</td>
                        <td nowrap>${formatDate(e.getDate())}</td>
                        <td>
                            <a href="/events/${e.getId()}"><i class="icon icon-screenshot"></i></a>
                            <a href="/events/${e.getId()}" target="_blank">${e.getName()}</a>
                        </td>
                        <td>${e.getLastNoteCount()}</td>
                        <td>${e.getActorCount()}</td>
                    </tr>`) as HTMLTableRowElement;

                    const liveDate = convertToJsDate(e.getDate(), 15); // next day 0:00 JST.
                    if (liveDate > new Date()) {
                        trDom.classList.add('warning');
                    }

                    tbodyDom.appendChild(trDom);
                    loadedCount += 1;
                });

                loadedIndicatorDom.innerText = loadedCount.toString();

                if (response.getHasNext()) {
                    enableLoadMoreButton();
                } else {
                    totalIndicatorDom.innerText = loadedCount.toString();
                }
            });
        }

        loadMoreButtonDom.addEventListener('click', loadMore);
        enableLoadMoreButton();

        if (autoLoadFirstPage) {
            loadMore();
        }
    }

    function placePage(placeId: string) {
        const tdDom = htmlToElement(`
        <div>
            ${header}
            <h3>Top Event List</h3>
        </div>`) as HTMLDivElement;

        const placeDetailsTableDom = document.getElementsByClassName('gb_place_detail_table')[0] || document.getElementsByClassName('mod_places_detail')[0];
        placeDetailsTableDom.parentNode.insertBefore(tdDom, placeDetailsTableDom.nextSibling);

        const request = new pb.RenderPlaceRequest().setId(placeId);
        treeDiagramService.renderPlace(request, null).then(response => {
            createEventList(new pb.QueryEventsRequest.EventFilter().setPlaceId(placeId), response.getKnownEventCount(), tdDom, false);
        });
    }

    function actorSnapshotsToDataPoints(snapshots: pb.RenderActorsResponse.ResponseItem.Snapshot[]): [number, number][] {
        if (snapshots.length === 0) {
            return [];
        }

        // JST 6am.
        const dates = snapshots.map(s => convertToJsDate(s.getDate(), -3));

        const dataPoints: [number, number][] = [];
        const copy = snapshots.slice();
        let lastFavoriteCount = 0;

        for (let date = dates[0]; date < new Date(); date.setDate(date.getDate() + 1)) {
            if (copy.length > 0 && date >= dates[0]) {
                lastFavoriteCount = copy[0].getFavoriteCount();
                dates.shift();
                copy.shift();
            }
            // Using array instead of object to make Highcharts Turbo mode happy.
            dataPoints.push([date.getTime(), lastFavoriteCount]);
        }
        return dataPoints;
    }

    function renderActorPage(actorId: string, renderChartTo: string, eventListDom?: HTMLElement) {
        const request = new pb.RenderActorsRequest().addId(actorId);
        treeDiagramService.renderActors(request, null).then(response => {
            const data = response.getItemsMap().get(actorId);

            if (eventListDom) {
                createEventList(new pb.QueryEventsRequest.EventFilter().addActorIds(actorId), data.getKnownEventCount(), eventListDom, false);
            }

            Highcharts.chart(renderChartTo, {
                chart: {zoomType: 'x'},
                xAxis: {type: 'datetime'},
                yAxis: {title: {text: undefined}},
                legend: {enabled: false},
                tooltip: {xDateFormat: '%Y-%m-%d'},
                series: [{
                    name: 'FavoriteCount',
                    type: 'areaspline',
                    data: actorSnapshotsToDataPoints(data.getSnapshotsList()),
                }],
            } as Options);
        });
    }

    function actorPage(actorId: string) {
        const tdDom = htmlToElement(`
        <div>
            ${header}
            <h3>Top Event List</h3>
        </div>`) as HTMLDivElement;

        const actorTitleDom = document.getElementsByClassName('gb_actors_title')[0] || document.getElementsByClassName('gb_blur_title')[0];
        actorTitleDom.parentNode.insertBefore(tdDom, actorTitleDom.nextSibling);

        const graphDom = htmlToElement(`
        <div>
            <div id="td_chart"></div>
        </div>`);
        const favoriteUsersDom = document.getElementsByClassName('gb_users_icon')[0] || document.getElementsByClassName('gb_listusericon')[0];
        favoriteUsersDom.parentNode.insertBefore(graphDom, favoriteUsersDom);

        renderActorPage(actorId, 'td_chart', tdDom);
    }

    function actorEventsPage(actorId: string) {
        const graphDom = htmlToElement(`
        <div>
            <div id="td_chart"></div>
        </div>`);

        const sidebarContentDom = document.getElementsByClassName('gb_ad_lrec')[0];
        if (sidebarContentDom) {
            sidebarContentDom.parentNode.insertBefore(graphDom, sidebarContentDom);
        } else {
            document.getElementsByClassName('mod_page')[0].appendChild(graphDom);
        }

        renderActorPage(actorId, 'td_chart');
    }

    function userPage() {
        const favoriteActorsDom = document.getElementsByClassName('gb_actors_list')[0] || document.getElementsByClassName('favorite_actor')[0];
        if (favoriteActorsDom) {
            const actorNames = new Map<string, string>();

            const actorDoms = favoriteActorsDom.getElementsByTagName('li');
            for (let i = 0; i < actorDoms.length; i++) {
                let count = actorDoms[i].className.match(/c(\d+)/)[1];
                const aEl = actorDoms[i].getElementsByTagName('a')[0];
                const splited = aEl.href.split('/');
                const actorId = splited[splited.length - 1];
                actorNames.set(actorId, aEl.textContent);
                aEl.textContent += ` (${count})`;
            }

            const containerForCompareActors = document.getElementsByClassName('span8')[0] || document.getElementsByClassName('mod_page')[0];
            if (actorNames.size <= renderActorsMax && containerForCompareActors) {
                containerForCompareActors.appendChild(htmlToElement(`
                    <div>
                        <button class="btn btn-block" type="button" id="td_user_compare_actors">Compare Favorites</button>
                        <div id="td_user_actor_favorites_chart_container"></div>
                    </div>`));

                const buttonEl = document.getElementById("td_user_compare_actors");
                buttonEl.addEventListener('click', () => {
                    compareActors(document.getElementById("td_user_actor_favorites_chart_container"), actorNames, () => buttonEl.remove());
                });
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

        const actorNames = new Map<string, string>();
        const selectedActorsDom = document.getElementById('td_selected_actors');
        const searchActorResultDom = document.getElementById('td_search_actor_result');
        const actorFavoritesChartContainerDom = document.getElementById('td_actor_favorites_chart_container');
        const eventListContainerDom = document.getElementById('td_event_list_container');

        const inputDom = document.getElementById('td_search_actor_input') as HTMLInputElement;
        inputDom.addEventListener('keyup', () => {
            searchActor(inputDom.value);
        });

        document.getElementById('td_compare_actors').addEventListener('click', () =>
            compareActors(actorFavoritesChartContainerDom, actorNames));
        document.getElementById('td_run_query').addEventListener('click', runQuery);

        function addActor(id: string, name: string) {
            if (actorNames.has(id)) return;
            actorNames.set(id, name);

            const buttonDom = htmlToElement(`
                <button class="btn" type="button"><i class="icon-minus"></i> ${name}</button>`) as HTMLButtonElement;
            selectedActorsDom.appendChild(buttonDom);

            buttonDom.addEventListener('click', () => {
                removeActor(id, buttonDom);
            });
        }

        function removeActor(id: string, el: HTMLElement) {
            if (!actorNames.has(id)) return;
            actorNames.delete(id);
            el.remove();
        }

        function searchActor(keyword: string) {
            if (!keyword) return;
            fetch(`/api/actors/search?&simple=1&limit=20&keyword=${keyword}`)
                .then(response => response.json())
                .then(data => {
                    if (inputDom.value !== keyword) return;
                    while (searchActorResultDom.firstChild) {
                        searchActorResultDom.removeChild(searchActorResultDom.firstChild);
                    }
                    if (data.results) {
                        data.results.forEach((item: { name: string, id: number }) => {
                            const itemDom = htmlToElement(`
                                <button class="btn" type="button"><i class="icon-plus"></i> ${item.name}</button>`);
                            itemDom.addEventListener('click', () => {
                                addActor(item.id.toString(), item.name);
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
            const filter = new pb.QueryEventsRequest.EventFilter()
            actorNames.forEach((_, a) => filter.addActorIds(a));
            createEventList(filter, undefined, eventListContainerDom, true);
        }
    }

    function compareActors(containerDom: HTMLElement, actorNames: Map<string, string>, callback?: () => void, extraVerticalLines?: { label: string, time: Date }[]) {
        while (containerDom.firstChild) {
            containerDom.removeChild(containerDom.firstChild);
        }

        containerDom.appendChild(htmlToElement(`<div id="td_chart_compare_actors"></div>`));

        const plotLines = extraVerticalLines.map(e => {
            return {value: e.time.getTime(), dashStyle: 'Dash', label: {text: e.label},}
        });

        const request = new pb.RenderActorsRequest();
        actorNames.forEach((_, i) => request.addId(i));

        treeDiagramService.renderActors(request, null).then(response => {
            const series: { name: string, type: 'spline', data: [number, number][] }[] = [];
            for (const [actorId, res] of response.getItemsMap().entries()) {
                series.push({
                    name: `${actorNames.get(actorId)}`,
                    type: 'spline',
                    data: actorSnapshotsToDataPoints(res.getSnapshotsList()),
                });
            }
            const lastValue = (l: [number, number][]) => l.length > 0 ? l[l.length - 1][1] : 0;
            series.sort((a, b) => lastValue(a.data) - lastValue(b.data));
            series.reverse();

            Highcharts.chart('td_chart_compare_actors', {
                chart: {zoomType: 'x'},
                xAxis: {
                    type: 'datetime',
                    plotLines: plotLines,
                },
                yAxis: {title: {text: undefined}},
                series: series,
                tooltip: {xDateFormat: '%Y-%m-%d', shared: true},
            } as Options);

            if (callback) callback();
        });
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

        const actorEventsPageRegex = /^https:\/\/www.eventernote.com\/actors\/(?:.+\/)?(\d+)\/?(events\/?(\?.*)?)?$/;
        const actorEventsPageMatch = url.match(actorEventsPageRegex);
        if (actorEventsPageMatch) {
            return actorEventsPage(actorEventsPageMatch[1]);
        }

        const userPageRegex = /^https:\/\/www.eventernote.com\/users(?:\/.+)?$/;
        const userPageMatch = url.match(userPageRegex);
        if (userPageMatch) {
            return userPage();
        }
    }

    main();
})();
