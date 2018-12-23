// ==UserScript==
// @name         TreeDiagram
// @namespace    https://www.sshz.org/
// @version      0.1
// @description  Make EventerNote Great Again
// @author       SSHZ.ORG
// @match        https://www.eventernote.com/*
// @require      https://cdnjs.cloudflare.com/ajax/libs/Chart.js/2.7.3/Chart.bundle.min.js
// ==/UserScript==

(function () {
    'use strict';

    function htmlToElement(html) {
        var template = document.createElement('template');
        html = html.trim();
        template.innerHTML = html;
        return template.content.firstChild;
    }

    let url = document.URL;

    let eventPageRegex = /https:\/\/www.eventernote.com\/events\/(\d+)/g;
    let eventPageMatch = eventPageRegex.exec(url);
    if (eventPageMatch !== null) {
        let canvasDom = htmlToElement('<canvas id="td_chart"></canvas>');

        let entryAreaDom = document.getElementById('entry_area');
        entryAreaDom.parentNode.insertBefore(canvasDom, entryAreaDom.nextSibling);

        let eventId = eventPageMatch[1];

        let xhr = new XMLHttpRequest();
        xhr.addEventListener("load", function (response) {
            let data = JSON.parse(xhr.responseText);
            let ctx = document.getElementById("td_chart");

            let tdChart = new Chart(ctx, {
                type: 'line',
                data: {
                    datasets: [{
                        label: 'NoteCount',
                        data: data.map(function (i) {
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

        xhr.open("GET", "https://treediagram.sshz.org/api/getNoteCountHistory?id=" + eventId, true);
        xhr.send();
    }
})();
