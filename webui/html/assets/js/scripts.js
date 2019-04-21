(function ($) {
    // "use strict";

    // 终端配置
    (function () {
        var terminal_jq = $('#terminal');
        var termical_height = $('#terminal').parent().height();
        var terminal = document.getElementById('terminal');
        // ws;

        (function () {
            ws = new WebSocket("ws:\/\/127.0.0.1:65500\/log");
            ws.onopen = function (evt) {
                print("Websocket OPEN\n");
            };
            ws.onclose = function (evt) {
                print("Websocket CLOSE\n");
                ws = null;
            };
            ws.onmessage = function (evt) {
                print(evt.data);
            };
            ws.onerror = function (evt) {
                print("Websocket ERROR: " + evt.data + "\n");
            };
            return ws;
        })();

        var auto_scroll = true;
        $('#auto_bottom').click(function () {
            terminal_jq.scrollTop(terminal_jq[0].scrollHeight);
            auto_scroll = true;
        });
        $('#stop_bottom').click(function () {
            auto_scroll = false;
        });

        var print = function (message) {
            var len = terminal_jq[0].scrollHeight - terminal_jq.scrollTop();
            if (len > termical_height + 100)
                auto_scroll = false;
            else if (len < termical_height + 10)
                auto_scroll = true;

            terminal.innerHTML = terminal.innerHTML.substr(-1024 * 100) + message;
            if (auto_scroll)
                terminal_jq.scrollTop(terminal_jq[0].scrollHeight);
        };
    })();


    // $('#terminal').scrollTop($('#terminal')[0].scrollHeight);

    /*================================
    Preloader
    ==================================*/
    // var preloader = $('#preloader');
    // $(window).on('load', function() {
    //     preloader.fadeOut('slow', function() { $(this).remove(); });
    // });

    //
    setInterval(function () {
        var $ = jQuery;
        var preloader = $('#preloader');
        preloader.fadeOut('slow', function () {
            $(preloader).remove();
        });
    }, 500);

    // 交易详情tab切换
    var txin = $('#txin');
    var txout = $('#txout');
    $('#txin_btn').click(function () {
        txout.removeClass("fade show active");
        txin.addClass("fade show active");
    });
    $('#txout_btn').click(function () {
        txout.addClass("fade show active");
        txin.removeClass("fade show active");
    });


    /*================================
    Start Footer resizer
    ==================================*/
    var e = function () {
        var e = (window.innerHeight > 0 ? window.innerHeight : this.screen.height) - 5;
        (e -= 67) < 1 && (e = 1), e > 67 && $(".main-content").css("min-height", e + "px")
    };
    $(window).ready(e), $(window).on("resize", e);

    /*================================
    sidebar menu
    ==================================*/
    // $("#menu").metisMenu();

    /*================================
    slimscroll activation
    ==================================*/
    $('.menu-inner').slimScroll({
        height: 'auto'
    });
    $('.nofity-list').slimScroll({
        height: '340px'
    });
    $('.timeline-area').slimScroll({
        height: '400px'
    });

    $('#wallets').slimScroll({
        height: '200px'
    });

    $('#tx_list').slimScroll({
        height: '410px'
    });


})(jQuery);