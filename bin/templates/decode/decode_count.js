var book_id = geturl2();
session = getCookie('session');

function count_init() {
    let reqBody = {
        "book_id": book_id, //图书编码
    };

    let dat = {
        'url': 'http://' + window.location.hostname + ':8000/barcodecount/',
        'method': 'POST',
        'req_body': JSON.stringify(reqBody)
    };

    $.ajax({
        url: 'http://' + window.location.hostname + ':8080/api',
        type: 'post',
        data: JSON.stringify(dat),
        headers: {'X-Session-Id': session},
        statusCode: {
            500: function () {
                console.log(null, "Internal error");
            }
        },
        complete: function (xhr, textStatus) {
            if (xhr.status >= 400) {
                console.log(null, "Error of Signin");
            }
        }
    }).done(function (data, statusText, xhr) {
        if (xhr.status >= 400) {
            console.log(null, "Error of Signin");
            return;
        }

        //console.log(data);
        if (data !== null) {
            $("#listNone").hide();
            let obj = JSON.parse(data);
            // console.log(obj);
            //console.log("长度：" + obj.barcodeIds.length);

            if (obj.code === "1") {

                let barcode = obj.message;

                if (barcode == null) {
                    return;
                }
                if (barcode.length == 0) {
                    return;
                }

                $("#jhm_num").text(barcode.all_count); //激活码总数
                $("#ysj_jhmNum").text(barcode.putaway_count); //已上架
                $("#wsj_num").text(barcode.noput_count + "个"); //未上架
                showPercent(barcode.putaway_count, barcode.all_count); //展示数据百分比
                $("#decode_num1").text(barcode.all_count); //激活码总数
                $("#decode_num2").text(barcode.use_count); //已激活解码锁
                if (barcode.all_count !== 0) {
                    $("#bar_left").css("height", "226px");
                }
                let _height = (barcode.use_count / barcode.all_count) * 226;
                $("#bar_mid").css("height", _height + "px"); //第二个柱子高度赋值
            } else if (obj.code === "2") {
                window.location.href = "/";
            }
        }
    });

    $(".icon_back").click(function () {//返回
        window.location.href = "/bookmanager";
    });
}

count_init(); //初始化方法
//进度方法计算
function showPercent(onNum, totalNum) {
    if (onNum !== 0) {
        let deg = (onNum / totalNum) * 360;
        if (Number(deg) <= 180) {
            $(".rightcircle").css("transform", "rotate(" + deg + "deg)");
        } else {
            $(".rightcircle").css("transform", "rotate(180deg)");
            $(".leftcircle").css("transform", "rotate(" + (deg - 180) + "deg)");
        }
    }

}