// ==UserScript==
// @name         TreeDiagram
// @namespace    https://www.sshz.org/
// @version      0.1.3
// @description  Make EventerNote Great Again
// @author       SSHZ.ORG
// @match        https://www.eventernote.com/*
// @require      https://cdnjs.cloudflare.com/ajax/libs/Chart.js/2.7.3/Chart.bundle.min.js
// ==/UserScript==

(function () {
    'use strict';

    function htmlToElement(html) {
        const template = document.createElement('template');
        html = html.trim();
        template.innerHTML = html;
        return template.content.firstChild;
    }

    function eventPage(eventId) {
        const tdDom = htmlToElement(`
            <div>
                <h2><ruby>樹形図の設計者<rt>ツリーダイアグラム</rt></ruby></h2>
                <canvas id="td_chart"></canvas>
                <div class="ma10 alert alert-info">
                    <span class="label">Total</span>
                    <span id="td_place_stats_total"></span>
                    <span class="label">Finished</span>
                    <span id="td_place_stats_finished"></span>
                </div>
            </div>`);

        const entryAreaDom = document.getElementById('entry_area');
        entryAreaDom.parentNode.insertBefore(tdDom, entryAreaDom.nextSibling);

        fetch(`https://treediagram.sshz.org/api/renderEvent?id=${eventId}`).then(function (response) {
            return response.json();
        }).then(function (data) {
            const totalStatsSpan = document.getElementById("td_place_stats_total");
            totalStatsSpan.innerHTML = `${data.place_stats_total.rank}/${data.place_stats_total.total}`;
            const finishedStatsSpan = document.getElementById("td_place_stats_finished");
            finishedStatsSpan.innerHTML = `${data.place_stats_finished.rank}/${data.place_stats_finished.total}`;

            const ctx = document.getElementById("td_chart");
            const tdChart = new Chart(ctx, {
                type: 'line',
                data: {
                    datasets: [{
                        label: 'NoteCount',
                        data: data.snapshots.map(function (i) {
                            return {
                                x: new Date(i.timestamp),
                                y: i.note_count,
                            };
                        }),
                        cubicInterpolationMode: 'monotone',
                        backgroundColor: 'rgba(54, 162, 235, 0.5)',
                        borderColor: 'rgb(54, 162, 235)',
                        borderWidth: 1
                    }]
                },
                options: {
                    scales: {
                        xAxes: [{
                            type: 'time'
                        }]
                    },
                    legend: {
                        display: false
                    }
                }
            });
        });
    }

    function placePage(placeId) {
        const tdDom = htmlToElement(`
            <div>
                <h2><ruby>樹形図の設計者<rt>ツリーダイアグラム</rt></ruby></h2>
                <h3>Top Events</h3>
                <table class="table table-striped s">
                    <tbody id="td_top_events_tbody"></tbody>
                </table>
            </div>`);

        const placeDetailsTableDom = document.getElementsByClassName('gb_place_detail_table')[0];
        placeDetailsTableDom.parentNode.insertBefore(tdDom, placeDetailsTableDom.nextSibling);

        const topEventsTbody = document.getElementById('td_top_events_tbody');

        fetch(`https://treediagram.sshz.org/api/renderPlace?id=${placeId}`).then(function (response) {
            return response.json();
        }).then(function (data) {

            data.top_events.forEach(function (e) {
                const trDom = htmlToElement(
                    `<tr>
                        <td>${e.date}</td>
                        <td><a href="/events/${e.id}">${e.name}</a></td>
                        <td>${e.last_note_count}</td>
                    </tr>`);

                topEventsTbody.appendChild(trDom);
            });
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
    }

    main();
})();
