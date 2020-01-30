// ==UserScript==
// @name         Better l-tike
// @namespace    https://www.sshz.org/
// @version      0.0.0.1
// @description  Let Me Paste ローソンさん
// @author       SSHZ.ORG
// @match        https://l-tike.com/*
// ==/UserScript==

(function () {
    'use strict';

    function main() {
        for (let input of document.getElementsByTagName('input')) {
            input.onpaste = undefined;
        }
    }

    main();
})();
