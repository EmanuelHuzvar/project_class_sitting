const loginIPSchool  = "http://10.2.130.45:9090/login" ;
const loginIPHome  = "http://192.168.1.10:9090/login" ;

document.getElementById('loginForm').addEventListener('submit', (event) => {
    event.preventDefault();
    const username = document.getElementById("username").value
    const password = document.getElementById("password").value

    const data = {
        username: username,
        password: password
    }

    const http = new XMLHttpRequest()
    http.onreadystatechange = function() {
        if (this.readyState == 4) {
            const json = JSON.parse(http.response)
            if(this.status == 200){
                sessionStorage.setItem("token", json.token)
                sessionStorage.setItem("username", json.username)
                sessionStorage.setItem("role", json.role)

                if(json.role == "user"){
                    window.location.assign("/course.html")
                } else {
                    window.location.assign("/adminCourse.html")
                }
            } else if (this.status == 401){
                document.getElementById("errorDiv").innerHTML = json.error
            }
        }
    };

    http.open("POST", loginIPSchool, true)
    http.send(JSON.stringify(data))
})