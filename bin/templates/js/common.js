// var urlstring = 'http://47.104.19.254:8080/south_sun/services/';
// //var urlstring = "http://192.168.18.6:8080/south_sun/services/";
// $.postajax = function (url1, data1, fun, flag, beforeFun, comleteFun) { //flag如果传值则表示同步请求，不传值则表示默认异步请求,beforeFun为发送ajax过程中的事件，comleteFun为请求完成后的事件
//     var _flag = true;
//     if (flag) {
//         _flag = false;
//     }
//     $.ajax({
//         type: "post",
//         url: urlstring + url1,
//         async: _flag,
//         beforeSend: beforeFun,
//         data: JSON.stringify(data1),
//         success: fun,
//         complete: comleteFun,
//         error: function (data) {
//             console.log(data);
//         }
//     });
// };

function geturl() { //获取页面之间的传值
    var reg = new RegExp("bookbarcode/(\\d+)", "i");
    var url = window.location.href;
    var r = url.match(reg);
    if (r != null) return unescape(r[1]);
    return null;
}

function geturl2() { //获取页面之间的传值
    var reg = new RegExp("bookbarcodecount/(\\d+)", "i");
    var url = window.location.href;
    var r = url.match(reg);
    if (r != null) return unescape(r[1]);
    return null;
}
//
// $("#header").load("../html/header.html", function () {
//     var str = '<div id="change_password" class="decode_div shadow">' +
//         '<div class="dialog_top">修改密码</div>' +
//         '<div class="decode_main">' +
//         '<div class="password_div clearfloat">' +
//         '<p>原密码：</p>' +
//         '<input type="password" class="decode_input1 decode_input" placeholder="请输入原来的密码" id="old_psw" />' +
//         '</div>' +
//         '<div class="password_div clearfloat">' +
//         '<p>新密码：</p>' +
//         '<input type="password" class="decode_input2 decode_input" placeholder="请输入新的密码" id="new_psw" />' +
//         '</div>' +
//         '<div class="password_div clearfloat">' +
//         '<p>重复输入：</p>' +
//         '<input type="password" class="decode_input3 decode_input" placeholder="请输入重复输入新密码" id="re_psw" />' +
//         '</div>' +
//         '<div class="clearfloat btn">' +
//         '<input type="button" value="确定" class="dialog_btn sure_btn" />' +
//         '<input type="button" value="取消" class="dialog_btn cancel_btn" />' +
//         "</div>" +
//         "</div>" +
//         "</div>"
//     $('body').append(str);
//     //后台退出登录
//     $('.exit').click(function () {
//         delCookie('us'); //删除cookie
//         sessionStorage.removeItem('cusinf');
//         window.location.href = "../login/login.html";
//
//     });
    // $('.change_password').click(function () {
    //     //console.log('修改密码！');
    //     windowShow("480px", "300px", $("#change_password"), $("#change_password .cancel_btn"), false);
    // });
    // $("#change_password .sure_btn").click(function () {
    //     console.log('点击了修改密码中的确认按钮！');
    //     var psw1 = $('.decode_input1').val();
    //     var psw2 = $('.decode_input2').val();
    //     var psw3 = $('.decode_input3').val();
    //     if (isNull(psw1) || isNull(psw2) || isNull(psw3)) {
    //         //console.log('kong');
    //         alert('请把新旧密码输入完整');
    //     } else {
    //         //console.log('bukong');
    //         if (psw2 === psw3) {
    //             $.postajax('user/updatePass', {
    //                 old_password: psw1,
    //                 new_password: psw2,
    //                 id: JSON.parse(sessionStorage.cusinf).id
    //             }, function (data) {
    //                 if (data.status == 0) {
    //                     alert('修改密码失败，请重试！')
    //                 } else if (data.status == 1) {
    //                     alert('修改成功');
    //                     window.location.replace('../login/login.html');
    //                 } else if (data.status == 2) {
    //                     alert('原密码错误！');
    //                 }
    //             })
    //         } else {
    //             alert('新旧密码请输入一致');
    //         }
    //     }
    // })
// });
//$("#footer").load("../html/footer.html", function() {

//});
function isNull(val) {
    if (val != '' && val != undefined && val != null && val != 'undefined' && val != 'null') {
        return false;//不等于空
    } else {
        return true;//等于空
    }
}

function setCookie(cname, cvalue, exmin) {
    var d = new Date();
    d.setTime(d.getTime() + (exmin * 60 * 1000));
    var expires = "expires=" + d.toUTCString();
    document.cookie = cname + "=" + cvalue + ";" + expires + ";path=/";
}

function getCookie(cname) {
    var name = cname + "=";
    var ca = document.cookie.split(';');
    for (var i = 0; i < ca.length; i++) {
        var c = ca[i];
        while (c.charAt(0) == ' ') {
            c = c.substring(1);
        }
        if (c.indexOf(name) == 0) {
            return c.substring(name.length, c.length);
        }
    }
    return "";
}
//获取cookie
// function getCookie(name) {
//     var arr, reg = new RegExp("(^| )" + name + "=([^;]*)(;|$)"); //正则匹配
//     if (arr = document.cookie.match(reg)) {
//         return unescape(arr[2]);
//     } else {
//         return null;
//     }
// };

//删除cookie
function delCookie(name) {
    var exp = new Date();
    exp.setTime(exp.getTime() - 1);
    var cval = getCookie(name);
    if (cval != null) {
        document.cookie = name + "=" + cval + ";expires=" + exp.toGMTString() + "; path=/";
    }

};

//把传进来的毫秒数装换成字符串类型然后输出,如果不传或者time怎返回具体时间,date怎返回日期,all返回全部的日期时间
function changeTime1(shijian, date1) {
    var show = new Date(Number(shijian));
    var year = show.getFullYear();
    var month = (show.getMonth() + 1) < 10 ? '0' + (show.getMonth() + 1) : show.getMonth() + 1
    var day = show.getDate() < 10 ? '0' + show.getDate() : show.getDate();
    var hour = show.getHours() < 10 ? '0' + show.getHours() : show.getHours();
    var minutes = show.getMinutes() < 10 ? '0' + show.getMinutes() : show.getMinutes();
    var seconds = show.getSeconds() < 10 ? '0' + show.getSeconds() : show.getSeconds();
    var stringdata = hour + ':' + minutes + ':' + seconds + ' ';
    var stringdata2 = year + '-' + month + '-' + day + ' ' + hour + ':' + minutes + ':' + seconds;
    var stringdata3 = year + '-' + month + '-' + day;
    if (date1 == '' || date1 == undefined || date1 == null || date1 == 'time') {
        return stringdata
    } else if (date1 == 'all') {
        return stringdata2;
    } else if (date1 == 'date') {
        return stringdata3;
    } else {
        return '请检查你的第二个参数,O(∩_∩)O谢谢!';
    }
}

//确认弹窗
function aresure(dom, str, string1, btn) { //string1是你想在中间展示的提示内容!,imgurl是图片的地址如果传false为默认的感叹号,btn传"1"就是一个按钮,"2"为两个按钮
    //	if($("#" + dom).size() <= 0) { //判断弹窗不存在再添加
    //		$(str).appendTo($("body"));
    //	}
    $("body").find($("#" + dom)).remove();
    $(str).appendTo($("body"));
    //console.log(string1);
    if (string1 !== undefined && string1 !== false && string1 !== '') {
        $('.dialog_main p').text(string1);
    }
    if (btn === 2) {

    } else if (btn === 1) {
        $('.cancel_btn').hide();
    }
    var layer = layui.layer;
    layer.open({
        type: 1,
        title: '',
        closeBtn: false,
        area: [],
        shadeClose: false,
        //btn: ["确认", "取消"],
        //btnAlign: 'c',
        shift: 0, //弹出的动画!
        resize: false, //是否可是拉动右下角的按钮,改变大小,false不可以
        scrollbar: false,
        end: function () {
            $("#" + dom).hide();
        },
        content: $('#' + dom), //这里content是一个DOM，注意：最好该元素要存放在body最外层，否则可能被其它的相对元素所影响
    });
    //	$('.suretipstit b').click(function() {
    //		layer.closeAll();
    //		console.log(1)
    //	});
}

//这是关闭的函数1,第一个参数是点击左边的按钮传入的参数,第二个参数是右边按钮传入的参数
function closearesure(dom, fnl, fnr) {
    $('body').find('#' + dom + ' .sure_btn').off("click").on('click', function () {
        if (fnl()) {
            layer.closeAll();
            //		console.log('确定');
        }
        ;
    })
    $('body').find('#' + dom + ' .cancel_btn').off("click").on('click', function () {
        layer.closeAll();
        console.log('取消');
        if (fnr) {
            fnr();
        }
    })
};

//自动关闭弹窗2,
function sureBox(txt, img) {
    $('#msginf').remove()
    var str = '<div id="msginf">' +
        '<div class="msginf_img">' +
        '<img src="../img/cg.png" alt="" />' +
        '</div>' +
        '<div class="msginf_tips">' + txt + '</div>' +
        '</div>';
    if ($(".msginf").size() <= 0) { //判断弹窗不存在再添加
        $(str).appendTo($("body"));
    }
    ;
    if (img != undefined || img != '' || img != false) { //替换页面中的图片
        $('#msginf').find('img').attr('src', img)
    }
    ;
    layer.open({
        skin: 'demo-class',
        type: 1,
        isOutAnim: true,
        time: 1000,
        area: ["280px", "200px"], //弹窗的大小
        content: $('#msginf'),
        closeBtn: false, //右上角的关闭按钮
        title: false,
        //title: ['提示', 'font-size:24px;background:#FFF;color:#333;padding:24px;height:24px;line-height:24px'],
        btnAlign: 'c',
    });
    $('#msginf').parent().parent().css('border-radius', '20px'); //改变弹窗的圆角
    window.setTimeout(function () {
        $('#msginf').hide().css('display', 'none')
    }, 1000)
    $('.layui-layer-hui .layui-layer-content').css('padding', '0')
};

//第一个参数是你要显示的汉字,第二个是你要显示的图片的地址,都可以不填或者传false,则默认是原理来的文字,布局和js同时引入
function msg(string1, imgurl) {
    var str = "<div id='msg' class='msg_div'>" +
        "<div><p class='msg_top'>提示</p>" +
        "<div class='msg_main' id='msg_main'><img src='../img/zhuyi.png'/>" +
        "<p></p></div>" +
        "</div>" +
        "</div>";
    if ($("#msg").size() <= 0) { //判断弹窗不存在再添加
        $(str).appendTo($("body"));
    }
    if (string1 != undefined && string1 != false && string1 != '') {
        $('#msg_main p').text(string1);
    }
    if (imgurl != undefined && imgurl != false && imgurl != '') {
        $('#msg img').attr('src', imgurl);
    }
    var layer = layui.layer;
    layer.open({
        type: 1,
        title: '',
        closeBtn: false,
        area: [],
        shadeClose: false,
        //btn: ["确认", "取消"],
        //btnAlign: 'c',
        shift: 0, //弹出的动画!
        resize: false, //是否可是拉动右下角的按钮,改变大小,false不可以
        scrollbar: false,
        time: 1500,
        end: function () {
            $("#msg").hide();
        },
        content: $('#msg'), //这里content是一个DOM，注意：最好该元素要存放在body最外层，否则可能被其它的相对元素所影响
    });
};

//2017-10-18李鑫,弹窗js封装,简化了弹窗的引入
function windowShow(boxWidth, boxHeight, dom, domClose, btnEven) { //弹窗封装，这里只是相关的弹出事件，没有布局(boxWidth:弹窗宽度,boxHeight:弹窗高度,dom:显示的dom布局,domClose:右上角关闭按钮)
    //btnEven为数组,其中包括点击按钮,确定以及取消的回调函数;[['确定','取消'],确定回调,取消回调];若为自定义按钮直接传入[false];
    layui.use(['layer', 'laypage', 'element'], function () {
        var layer = layui.layer,
            laypage = layui.laypage,
            element = layui.element();
        domClose.click(function () {
            //layer.close(layer.index);
            layer.closeAll();
            dom.css('display', 'none')
        });
    });
    layer.open({
        skin: 'demo-class',
        type: 1,
        area: [boxWidth, boxHeight], //弹窗的大小
        content: dom,
        closeBtn: 0, //右上角的关闭按钮
        scrollbar: false,
        title: false,
        //title: ['提示', 'font-size:24px;background:#FFF;color:#333;padding:24px;height:24px;line-height:24px'],
        btnAlign: 'c',
        btn: btnEven[0],
        yes: function (index, layero) {
            //按钮【按钮一】的回调
            dom.hide().css('display', 'none'); //
            layer.close(index); //如果设定了yes回调，需进行手工关闭
            btnEven[1](); //确定回调函数
        },
        btn2: function (index, layero) {
            //按钮【按钮二】的回调
            dom.hide().css('display', 'none')
            layer.close(index); //如果设定了yes回调，需进行手工关闭
            //return false 开启该代码可禁止点击该按钮关闭
            btnEven[2](); //取消回调函数
        },
        cancel: function () {
            //右上角关闭回调
            dom.hide().css('display', 'none')
            //return false 开启该代码可禁止点击该按钮关闭
        }
    });
};

//windowShow关闭弹窗,这里可以自己写
function winsBtnClose(dom, btnSure, btnFalse) {
    if (btnSure != undefined && btnSure != '') {
        btnSure[0].click(function () {
            btnSure[1](); //确定函数的回调
        });
    }
    if (btnFalse != undefined && btnFalse != '') {
        btnFalse[0].click(function () {
            layer.closeAll();
            dom.css('display', 'none')
            btnFalse[1](); //确定函数的回调
        })
    }
};



