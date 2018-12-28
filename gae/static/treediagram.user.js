// ==UserScript==
// @name         TreeDiagram
// @namespace    https://www.sshz.org/
// @version      0.1.4
// @description  Make EventerNote Great Again
// @author       SSHZ.ORG
// @match        https://www.eventernote.com/*
// @require      https://cdnjs.cloudflare.com/ajax/libs/Chart.js/2.7.3/Chart.bundle.min.js
// ==/UserScript==

(function () {
    'use strict';

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

        fetch(`${serverBaseAddress}/api/renderEvent?id=${eventId}`).then(function (response) {
            return response.json();
        }).then(function (data) {
            const totalStatsSpan = document.getElementById("td_place_stats_total");
            totalStatsSpan.innerHTML = `${data.placeStatsTotal.rank}/${data.placeStatsTotal.total}`;
            const finishedStatsSpan = document.getElementById("td_place_stats_finished");
            finishedStatsSpan.innerHTML = `${data.placeStatsFinished.rank}/${data.placeStatsFinished.total}`;

            const ctx = document.getElementById("td_chart");
            const tdChart = new Chart(ctx, {
                type: 'line',
                data: {
                    datasets: [{
                        label: 'NoteCount',
                        data: data.snapshots.map(function (i) {
                            return {
                                x: new Date(i.timestamp),
                                y: i.noteCount,
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

    function createEventList(argString, totalCount, domToAppend) {
        let currentPage = 1; // 1-based for display.
        const totalPage = Math.floor((totalCount - 1) / 10) + 1;

        const eventListDom = htmlToElement(`
            <div>
                <table class="table table-hover s" style="margin-bottom: 0">
                    <tbody id="td_event_list_tbody"></tbody>
                </table>
                <div class="pagination pagination-centered">
                    <ul>
                        <li id="td_event_list_previous_page_button"><a href="#">&lt;</a></li>
                        <li class="active"><span>
                            <span id="td_event_list_page_indicator"></span> / ${totalPage} (${totalCount})
                        </span></li>
                        <li id="td_event_list_next_page_button"><a href="#">&gt;</a></li>
                    </ul>
                </div>
            </div>`);
        domToAppend.appendChild(eventListDom);

        const pageIndicatorDom = document.getElementById('td_event_list_page_indicator');

        function goToPage(newPage) {
            const tbodyDom = document.getElementById('td_event_list_tbody');
            while (tbodyDom.firstChild) {
                tbodyDom.removeChild(tbodyDom.firstChild);
            }

            fetch(`${serverBaseAddress}/api/queryEvents?page=${newPage}&${argString}`).then(function (response) {
                return response.json();
            }).then(function (data) {
                data.forEach(function (e) {
                    const trDom = htmlToElement(
                        `<tr>
                        <td>${e.date}</td>
                        <td><a href="/events/${e.id}">${e.name}</a></td>
                        <td>${e.lastNoteCount}</td>
                    </tr>`);

                    tbodyDom.appendChild(trDom);
                });

                currentPage = newPage;
                pageIndicatorDom.innerText = currentPage;
            });
        }

        const previousPageButton = document.getElementById('td_event_list_previous_page_button');
        previousPageButton.addEventListener('click', function () {
            if (currentPage - 1 < 1) {
                return;
            }
            goToPage(currentPage - 1);
        });

        const nextPageButton = document.getElementById('td_event_list_next_page_button');
        nextPageButton.addEventListener('click', function () {
            if (currentPage + 1 > totalPage) {
                return;
            }
            goToPage(currentPage + 1);
        });

        goToPage(1);
    }

    function placePage(placeId) {
        const tdDom = htmlToElement(`
            <div>
                <h2><ruby>樹形図の設計者<rt>ツリーダイアグラム</rt></ruby></h2>
                <h3>Event List</h3>
            </div>`);

        const placeDetailsTableDom = document.getElementsByClassName('gb_place_detail_table')[0];
        placeDetailsTableDom.parentNode.insertBefore(tdDom, placeDetailsTableDom.nextSibling);

        fetch(`${serverBaseAddress}/api/renderPlace?id=${placeId}`).then(function (response) {
            return response.json();
        }).then(function (data) {
            createEventList(`place=${placeId}`, data.knownEventCount, tdDom);
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
