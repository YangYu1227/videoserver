let page = 1;
session = getCookie('session');

function list_init() {
    //console.log(getCookie('session'));

    let dat = {
        'url': 'http://' + window.location.hostname + ':8000/bookmanagerlist',
        'method': 'GET',
        'req_body': ''
    };

    $.ajax({
        url: 'http://' + window.location.hostname + ':8080/api',
        type: 'POST',
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
        $("#table_list").html("");
        if (data !== null) {
            //console.log(data);
            var obj = JSON.parse(data);

            if (obj.code === "1") {
                $("#listNone").hide();

                let msg = obj.message;
                let books = msg.books;

                $.each(books, function (index, item) {
                    //console.log(item);
                    var str = '<ul id=' + item.book_id + ' class="clearfloat"><li class="width120">' + item.id + '</li>';
                    str += '<li class="width280 book_name">' + item.book_name + '</li>';
                    str += '<li class="width280 book_num">' + item.book_id + '</li>';
                    str += '<li class="width220">' + item.barcode_count + '</li>';
                    str += '<li class="width300">';
                    str += '<button class="pointer bookBtn_edit edit">【编辑图书】</button>';
                    str += '<button class="pointer bookBtn_jsm">【解锁码】</button>';
                    str += '<button class="pointer bookBtn_tj">【统计】</button>';
                    str += '</li>';
                    str += '</li></ul>';
                    $("#table_list").append(str);
                });
                fenye(Math.ceil(data.total / 10), page, "pageDemo");
                clickAddButton(); //点击事件
            } else if (obj.code === "2") {
                window.location.href = "/";
            } else {
                $("#listNone").show();
            }
        } else {
            $("#listNone").show();
        }
    });
}


list_init(); //初始化数据

function clickAddButton() {
    //	新增图书
    $(".bookBtn_add").off("click").on("click", function () {
        var parentId = $(this).parent().parent().attr("id");
        var dom;
        dom = "book_add";
        var str = '<div id="' + dom + '" class="shadow book_edit">';
        str += '<div class="dialog_top">新增图书</div>';
        str += '<div class="book_editMain">';
        str += '<div class="book_editLine">';
        str += '<span>图书名称：</span><input type="text" class="book_input" id="book_name" placeholder="请输入图书名称" max-length="20" />';
        str += '<span class="book_msg">图书名称不能为空</span></div>';
        str += '<div class="book_editLine">';
        str += '<span>图书编号：</span><input type="text" class="book_input" id="book_id" placeholder="请输入图书编号" max-length="20" />';
        str += '<span class="book_msg">图书编号不能为空</span></div>';
        str += '<div class="clearfloat btn_div">';
        str += '<input type="button" value="确定" class="dialog_btn sure_btn" />';
        str += '<input type="button" value="取消" class="dialog_btn cancel_btn" />';
        str += '</div></div></div>';
        aresure(dom, str, "", 2);
        close_dialog(dom);
    });

    // 编辑图书
    $(".bookBtn_edit").off("click").on("click", function () {
        var parentId = $(this).parent().parent().attr("id");
        var dom;
        var _name = $(this).parent().siblings(".book_name").text();
        var _num = $(this).parent().siblings(".book_num").text();
        dom = "book_edit";
        var str = '<div id="' + dom + '" class="shadow book_edit" data-id="' + parentId + '">';
        str += '<div class="dialog_top">编辑图书</div>';
        str += '<div class="book_editMain">';
        str += '<div class="book_editLine">';
        str += '<span>图书名称：</span><input type="text" class="book_input" id="book_name" value="' + _name + '" maxlength="20" />';
        str += '<span class="book_msg">图书名称不能为空</span></div>';
        str += '<div class="book_editLine">';
        str += '<span>图书编号：</span><input type="text" class="book_input" id="book_id" value="' + _num + '" maxlength="20" />';
        str += '<span class="book_msg">图书编号不能为空</span></div>';
        str += '<div class="clearfloat btn_div">';
        str += '<input type="button" value="确定" class="dialog_btn sure_btn" />';
        str += '<input type="button" value="取消" class="dialog_btn cancel_btn" />';
        str += '</div></div></div>';
        aresure(dom, str, "", 2);
        close_dialog(dom);
    });

    //解码锁
    $(".bookBtn_jsm").click(function () {
        var book_id = $(this).parent().parent().attr("id");
        window.location.href = "/bookbarcode/" + book_id;
    });
    //统计
    $(".bookBtn_tj").click(function () {
        var book_id = $(this).parent().parent().attr("id");
        window.location.href = "/bookbarcodecount/" + book_id;
    });

    //主页
    $(".bookBtn_Home").click(function () {
        window.location.href = "/bookmanager";
    });

    //修改密码
    $(".bookBtn_Change").click(function () {
        windowShow("480px", "324px", $("#change_password"), $("#change_password .cancel_btn"), false);
        $('#change_password .old_psw').val("");
        $('#change_password .new_psw').val("");
        $('#change_password .re_psw').val("");
    });

    $("#change_password .sure_btn").click(function () {
        // console.log('点击了修改密码中的确认按钮！');
        var old_psw = $('.old_psw').val();
        var new_psw = $('.new_psw').val();
        var re_psw = $('.re_psw').val();

        // console.log(old_psw);
        // console.log(new_psw);
        // console.log(re_psw);


        if (isNull(old_psw) || isNull(new_psw) || isNull(re_psw)) {
            $(".decode_msg").text("请把新旧密码输入完整").css("visibility", "visible");
        } else {
            //console.log('bukong');
            if (new_psw === re_psw) {
                // console.log("POST");
                $(".decode_msg").css("visibility", "hidden");


                let reqBody = {
                    "user_name": "admin", //用户名
                    "old_pwd": old_psw, //旧的密码
                    "new_pwd": new_psw, //新的密码
                };

                let dat = {
                    'url': 'http://' + window.location.hostname + ':8000/changepassword/',
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

                    console.log(data);
                    if (data !== null) {
                        let obj = JSON.parse(data);

                        if (obj.code === "1") {
                            $("#change_password .cancel_btn").click();
                            windowShow("480px", "184px", $("#change_password_msg"), $("#change_password_msg .cancel_btn"), false);
                            $("#change_password_msg .decode_msg").text(obj.message).css("visibility", "visible");
                        } else if (obj.code === "2") {
                            window.location.href = "/";
                        }else {
                            $("#change_password .cancel_btn").click();
                            windowShow("480px", "184px", $("#change_password_msg"), $("#change_password_msg .cancel_btn"), false);
                            $("#change_password_msg .decode_msg").text(obj.message).css("visibility", "visible");
                        }

                    } else {
                        $("#change_password .cancel_btn").click();
                        windowShow("480px", "184px", $("#change_password_msg"), $("#change_password_msg .cancel_btn"), false);
                        $("#change_password_msg .decode_msg").text("修改失败").css("visibility", "visible");
                    }
                });
            } else {
                $(".decode_msg").text("新旧密码请输入一致").css("visibility", "visible");
            }
        }
    });
}

//新增修改图书点击确定取消按钮
function close_dialog(dom) {
    //新增修改图书点击确定取消按钮
    $("#" + dom + " .sure_btn").off("click").on('click', function () {
        var _parent = $(this).parents(".book_edit").attr("id");
        var dat;
        var reqBody;
        if (_parent === "book_edit") { //编辑图书
            reqBody = {
                "id": $("#book_edit").attr("data-id"),
                "book_name": $("#book_edit #book_name").val(), //图书名称
                "book_id": $("#book_edit #book_id").val() //图书编号
            };
            console.log(reqBody);
            dat = {
                'url': 'http://' + window.location.hostname + ':8000/updatebook/',
                'method': 'POST',
                'req_body': JSON.stringify(reqBody)
            };
        } else { //新增图书
            reqBody = {
                "book_name": $("#book_add #book_name").val(), //图书名称
                "book_id": $("#book_add #book_id").val() //图书编号
            };
            dat = {
                'url': 'http://' + window.location.hostname + ':8000/addbook/',
                'method': 'POST',
                'req_body': JSON.stringify(reqBody)
            };
        }
        if (val_flag(dom)) {
            console.log(dat);
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
                if (data !== null) {
                    let obj = JSON.parse(data);

                    if (obj.code === "1") {
                        layer.closeAll();
                        page = 1;
                        list_init();
                    } else if (obj.code === "2") {
                        window.location.href = "/";
                    } else {
                        msg(obj.message);
                    }
                }
            });
        }
    });
    //新增修改图书点击确定取消按钮
    $("#" + dom + " .cancel_btn").off("click").on('click', function () {
        layer.closeAll();
    });
}

//校验是否为空
function val_flag(dom) {
    var _bookName = $("#" + dom + " #book_name").val(); //图书名称
    var _bookId = $("#" + dom + " #book_id").val(); //图书编号
    var _flag = false;
    if (!_bookName.trim()) {
        $("#book_name").siblings(".book_msg").css("visibility", "visible");
        _flag = false;
        return false;
    } else {
        $("#book_name").siblings(".book_msg").css("visibility", "hidden");
    }
    if (!_bookId.trim()) {
        $("#book_id").siblings(".book_msg").css("visibility", "visible");
        _flag = false;
        return false;
    } else {
        $("#book_name").siblings(".book_msg").css("visibility", "hidden");
    }
    if (_bookName.trim() && _bookId.trim()) {
        _flag = true;
    }
    return _flag;
}

//分页控件
function fenye(pages, currindex, pageId) {
    layui.use(['layer', 'laypage', 'element', 'laydate'], function () {
        var laypage = layui.laypage;
        laypage({
            cont: pageId, //分页容器的id
            pages: pages, //总页数
            skin: '#fff', //自定义选中色值
            skip: true, //开启跳页
            curr: page, //当前页码
            jump: function (obj, first) {
                if (!first) {
                    //layer.msg('第' + obj.curr + '页');
                    //					homeorder_postData.pages_number = (obj.curr - 1) * 10 + 1;
                    page = obj.curr;
                    list_init();
                }
            }
        });
    })
}