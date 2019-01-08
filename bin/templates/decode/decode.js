let decode_page = 1;
let book_id = geturl();
let decode_num = 0;//解锁码数量id
let barcodeId = "";
session = getCookie('session');

function decode_init() {
    data_init();
    click_fun(); //点击事件
}

decode_init(); //初始化数据
function data_init() {
    $("#decode_table").empty();
    decode_num = 0;
    $("#jsm_num").text(0);

    let reqBody = {
        "book_id": book_id, //图书编码
        // "barcode_id": $("#search_text").val(), //解锁码编号
        // "pages_number": decode_page,
        // "status": $(".stated").attr("data-state") //状态(0：全部；1：已上架；2：未上架)
    };

    let dat = {
        'url': 'http://' + window.location.hostname + ':8000/findbarcodelist/',
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
            //console.log(obj);
            //console.log("长度：" + obj.barcodeIds.length);

            if (obj.code === "1") {
                let msg = obj.message;
                //console.log(msg);

                if (msg.barcodeIds == null){
                    $("#listNone").show();
                    $("#jsm_num").text('0');
                    return;
                }

                if (msg.barcodeIds.length === 0) {
                    $("#listNone").show();
                    $("#jsm_num").text('0');
                    return;
                }

                $.each(msg.barcodeIds, function (index, item) {
                    //console.log("aaaaaa");
                    let str = '<li class="序号" style="width: 110px;">' + item.id + '</li>';
                    str += '<li class="解码锁编号" style="width: 310px;">' + item.barcode_id + '</li>';
                    str += '<li class="解锁码数量" style="width: 140px;">' + item.all_count + '</li>';
                    str += '<li class="已解锁数量" style="width: 140px;">' + item.use_count + '</li>';
                    str += '<li class="创建时间" style="width: 170px;">' + item.create_time + '</li>';
                    if (item.status === "0") {//状态：0，未上架；1：上架
                        str += '<li class="状态" style="width: 110px;">未上架</li>';
                        str += '<li class="操作" style="width: 220px;" barcode_id = "' + item.barcode_id + '" status = "' + item.status + '">';
                        str += '<button class="pointer bookBtn_sj">【上架】</button>';
                        str += '<button class="pointer bookBtn_del">【删除】</button>';
                    } else {
                        str += '<li class="状态" style="width: 110px;">已上架</li>';
                        str += '<li class="操作" style="width: 220px;" barcode_id = "' + item.barcode_id + '" status = "' + item.status + '">';
                        str += '<button class="pointer bookBtn_export">【导出】</button>';
                        str += '<button class="pointer bookBtn_del">【删除】</button>';
                    }
                    str += '</li></ul>';
                    $("#decode_table").append(str);
                    decode_num += item.all_count;
                    operation_button_click();
                });
                $("#jsm_num").text(decode_num);
                fenye(Math.ceil(msg.barcodeIds.length / 10)); //分页
            } else if (obj.code === "2") {
                window.location.href = "/";
            } else {
                $("#listNone").show();
                $("#jsm_num").text('0');
            }
        } else {
            $("#listNone").show();
            $("#jsm_num").text('0');
        }
    });
}

//操作下的两个按钮：上架，删除
function operation_button_click() {
    $(".bookBtn_sj").off("click").on("click", function () {
        console.log("上架");
        var barcode_id = $(this).parent().attr("barcode_id");
        //console.log(barcode_id);
        putaway(barcode_id);
    });

    $(".bookBtn_del").off("click").on("click", function () {
        console.log("删除");

        var status = $(this).parent().attr("status");

        barcodeId = $(this).parent().attr("barcode_id");

        if (status === "1") {
            windowShow("480px", "300px", $("#delete_div"), $("#delete_div .cancel_btn"), false);
        } else {
            windowShow("480px", "300px", $("#delete_div2"), $("#delete_div2 .cancel_btn"), false);
        }
    });

    // 导出解锁码按钮
    $(".bookBtn_export").click(function () {
        console.log("导出解锁码点击");
        var barcode_id = $(this).parent().attr("barcode_id");
        exportBarcode({
            url: 'http://' + window.location.hostname + ':8000/exportBarcode.xlsx', //请求的url
            data: {
                barcode_id,
                book_id,
            }
        })
    });
}

//上架
function putaway(barcode_id) {
    let reqBody = {
        "barcode_id": barcode_id, //解锁码编号
    };

    let dat = {
        'url': 'http://' + window.location.hostname + ':8000/changebarcodestatus/',
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
            let obj = JSON.parse(data);

            if (obj.code === "1") {
                data_init();
            } else if (obj.code === "2") {
                window.location.href = "/";
            }
        }
    });
}

// 删除
function deleteBarcode() {
    let reqBody = {
        "barcode_id": barcodeId, //解锁码编号
    };

    let dat = {
        'url': 'http://' + window.location.hostname + ':8000/deletebarcode/',
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
            let obj = JSON.parse(data);

            if (obj.code === "1") {
                data_init();
            } else if (obj.code === "2") {
                window.location.href = "/";
            }
        }
    });
}

// 导出
function exportBarcode(options) {

    var config = $.extend(true, {method: 'post'}, options);
    var $iframe = $('<iframe id="down-file-iframe" />');
    var $form = $('<form target="down-file-iframe" method="' + config.method + '" />');
    $form.attr('action', config.url);
    for (var key in config.data) {
        $form.append('<input type="hidden" name="' + key + '" value="' + config.data[key] + '" />');
    }
    $iframe.append($form);
    $(document.body).append($iframe);
    $form[0].submit();
    $iframe.remove();
}

function click_fun() {
    //主页
    $(".bookBtn_Home").click(function () {
        window.location.href = "/bookmanager";
    });

    //点击搜索按钮
    $("#search_btn").click(function () {
        decode_page = 1;
        data_init();
    });
    //新增解锁码按钮
    $(".decode_btn").click(function () {
        console.log("新增解锁码点击");
        windowShow("480px", "300px", $("#decode_div"), $("#decode_div .cancel_btn"), false);
    });

    //新增解锁码确定按钮
    $("#decode_div .sure_btn").click(function () {
        console.log("增加解锁码面板的确定按钮");
        let add_num = $("#decode_addText").val();

        if (0 < add_num && add_num <= 5000 && add_num !== '') {
            $(".decode_msg").css("visibility", "hidden");

            let reqBody = {
                "add_num": add_num, //增加解锁码的数量
                "book_id": book_id, //图书编码
            };

            let dat = {
                'url': 'http://' + window.location.hostname + ':8000/addbarcode/',
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
                    let obj = JSON.parse(data);

                    if (obj.code === "1") {
                        layer.closeAll();
                        $("#decode_div").css('display', 'none');
                        $("#decode_addText").val("");
                        decode_page = 1;
                        data_init();
                    } else if (obj.code === "2") {
                        window.location.href = "/";
                    }
                }
            });
        } else if (add_num.trim() > 5000) {
            $(".decode_msg").text("解锁码数量不能超过5000").css("visibility", "visible");
        } else if (add_num === '') {
            $(".decode_msg").text("解锁码数量不能为空").css("visibility", "visible");
        }

    });

    $("#delete_div .sure_btn").click(function () {
        console.log("删除已上架解锁码的确定按钮");
        let text = $("#input_Text").val();

        if (text === '') {
            $(".decode_msg").text("文本不能为空").css("visibility", "visible");
            return;
        }

        if (text !== '确认删除') {
            $(".decode_msg").text("文本错误").css("visibility", "visible");
            return;
        }

        $(".decode_msg").css("visibility", "hidden");
        $("#delete_div .cancel_btn").click();
        windowShow("480px", "300px", $("#delete_div2"), $("#delete_div2 .cancel_btn"), false);
    });

    $("#delete_div2 .sure_btn").click(function () {
        $(".decode_msg").css("visibility", "hidden");

        $("#delete_div2 .cancel_btn").click();
        deleteBarcode();
    });
}


//分页控件
function fenye(pages) {
    layui.use(['layer', 'laypage', 'element', 'laydate'], function () {
        let laypage = layui.laypage;
        laypage({
            cont: "decode_page", //分页容器的id
            pages: pages, //总页数
            skin: '#fff', //自定义选中色值
            skip: true, //开启跳页
            curr: decode_page, //当前页码
            jump: function (obj, first) {
                if (!first) {
                    decode_page = obj.curr;
                    data_init();
                }
            }
        });
    })
}

//选择框选择
function radio_fun() {
    //点击选择
    $(".radio_list").click(function () {
        let obj = $(this).find("img");
        let _flag = $(this).find("img").hasClass("checked");
        let _all = $(this).hasClass("radio_all");
        if (_all) { //全选
            if (_flag) { //如果选中了
                $(".radio_list").find("img").attr("src", "../img/check_box.png");
                $(".radio_all").find("img").attr("src", "../img/check_all.png");
                $(".radio_list").find("img").removeClass("checked");
            } else {
                $(".radio_list").find("img").attr("src", "../img/check_boxed.png");
                $(".radio_all").find("img").attr("src", "../img/check_alled.png");
                $(".radio_list").find("img").addClass("checked");
            }
        } else { //单个点击
            if (_flag) { //如果选中了
                if ($(".radio_all").find("img").hasClass("checked")) { //全选
                    $(".radio_all").find("img").removeClass("checked");
                    $(".radio_all").find("img").attr("src", "../img/check_all.png");
                }
                obj.attr("src", "../img/check_box.png");
                obj.removeClass("checked");
            } else {
                obj.attr("src", "../img/check_boxed.png");
                obj.addClass("checked");
            }

        }
    });
}