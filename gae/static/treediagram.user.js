// ==UserScript==
// @name         TreeDiagram
// @namespace    https://www.sshz.org/
// @version      0.1.2
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

    const url = document.URL;

    const eventPageRegex = /https:\/\/www.eventernote.com\/events\/(\d+)/g;
    const eventPageMatch = eventPageRegex.exec(url);
    if (eventPageMatch !== null) {
        const canvasDom = htmlToElement(`
            <div>
                <canvas id="td_chart"></canvas>
                <div class="ma10 alert alert-info">
                    <span class="label">Total</span>
                    <span id="td_place_stats_total"></span>
                    <span class="label">Finished</span>
                    <span id="td_place_stats_finished"></span>
                </div>
            </div>`);

        const entryAreaDom = document.getElementById('entry_area');
        entryAreaDom.parentNode.insertBefore(canvasDom, entryAreaDom.nextSibling);

        const eventId = eventPageMatch[1];

        const xhr = new XMLHttpRequest();
        xhr.addEventListener("load", function (response) {
            const data = JSON.parse(xhr.responseText);

            const totalStatsSpan = document.getElementById("td_place_stats_total");
            totalStatsSpan.innerHTML = data.place_stats_total.rank + "/" + data.place_stats_total.total;
            const finishedStatsSpan = document.getElementById("td_place_stats_finished");
            finishedStatsSpan.innerHTML = data.place_stats_finished.rank + "/" + data.place_stats_finished.total;

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

        xhr.open("GET", "https://treediagram.sshz.org/api/renderEvent?id=" + eventId, true);
        xhr.send();
    }
})();
