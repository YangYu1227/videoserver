function login_init() {
    DEFAULT_COOKIE_EXPIRE_TIME = 30;

    uname = '';
    session = '';

    session = getCookie('session');
    uname = getCookie('username');


    $("#login_btn").click(function () {
        var username = $("#login_user").val();
        var pwd = $("#login_password").val();

        if (username === '' || pwd === '') {
            callback(null, err);
        }

        var reqBody = {
            'user_name': username,
            'pwd': pwd
        };

        var dat = {
            'url': 'http://' + window.location.hostname + ':8000/login/',
            'method': 'POST',
            'req_body': JSON.stringify(reqBody)
        };

        $.ajax({
            url: 'http://' + window.location.hostname + ':8080/api',
            type: 'post',
            data: JSON.stringify(dat),
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

            let obj = JSON.parse(data);

            //console.log(obj);

            if (obj.code === "0") {
                let msg = obj.message;
                let login_msg = msg.login_msg;
                windowShow("480px", "300px", $("#login_msg"), $("#login_msg .cancel_btn"), false);
                $("#login_msg .decode_msg").text(login_msg).css("visibility", "visible");
            } else {
                let msg = obj.message;
                let session_id = msg.session_id;

                setCookie("session", session_id, DEFAULT_COOKIE_EXPIRE_TIME);
                setCookie("username", uname, DEFAULT_COOKIE_EXPIRE_TIME);

                let aaa = getCookie('session');
                //console.log(aaa);
                window.location.href = "/bookmanager";
            }
        });
    })
}

login_init(); //初始化数据