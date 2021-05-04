var http = require("http");

const userId = "124582164851";
const body = JSON.stringify({
  userId,
});

const options = {
  hostname: "localhost",
  port: "4000",
  path: "/users",
  method: "POST",
  headers: {
    "content-type": "application/json",
    authorization:
      "Bearer eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCIsImtpZCI6ImZaSGlDMjBlUl9PS2ZoZFJwRXQtRCJ9.eyJodHRwczovL2NvdWNoLXBvdGF0b2VzLmNvbS9jbGFpbXMvIjpbImFkbWluIl0sImdpdmVuX25hbWUiOiJYaHVsaW8iLCJmYW1pbHlfbmFtZSI6IkRvZGEiLCJuaWNrbmFtZSI6InhodWxpb2RvLmFrc2hpIiwibmFtZSI6IlhodWxpbyBEb2RhIiwicGljdHVyZSI6Imh0dHBzOi8vbGg0Lmdvb2dsZXVzZXJjb250ZW50LmNvbS8tSjVMY2tHQ0liR2svQUFBQUFBQUFBQUkvQUFBQUFBQUFBQUEvQU1adXVjbjZtTGozVDFoU3FRMFJWMTBrdWoyU0ItQy16QS9zOTYtYy9waG90by5qcGciLCJsb2NhbGUiOiJlbiIsInVwZGF0ZWRfYXQiOiIyMDIxLTAyLTIyVDEwOjE3OjEzLjM4OFoiLCJlbWFpbCI6InhodWxpb2RvLmFrc2hpQGdtYWlsLmNvbSIsImVtYWlsX3ZlcmlmaWVkIjp0cnVlLCJpc3MiOiJodHRwczovL2Rldi1wczVkcXFpcy5ldS5hdXRoMC5jb20vIiwic3ViIjoiZ29vZ2xlLW9hdXRoMnwxMDM4ODg3OTQ4NzcxMDI4MDgzNTciLCJhdWQiOiJEcWJDdlp0TDhjbjVwbERsYTlUWWxySkxoV0lYcFp0ViIsImlhdCI6MTYxMzk4OTAzOSwiZXhwIjoxOTI5Mzg4OTk5LCJub25jZSI6IlFVeEdaazFxVmpKQ1VsRm5lV2xCYldVMlZHSllZMjFxY241VFZHTk1abmRoYUV0eGNuSTNWbG8wYkE9PSJ9.VxsA24YQO2hu_qvOrff9Vv-WxRQ7ukgMOq5xCrcN53olBNk8KydFqYxvYbfv1_AoRgg90Hb_R7hPsnarzHYs3t2zU4cPtfJarre6M2_1azV4PZNYg1sRgKCe7muqIwS2r_zUAGekTdg70A6ZmBwzli5mkYvl-OU9KGLl8YpOfT5QNc0O-FjbXboin2KvrouRau7XulPXY5NPYiZV_MR7qn3zDmUuBQUYJBag0znHsy5S100-TEeDbyU7UcKZyMyfEtO7uxmx7XatOdYRMAyPd0QoR7LJ8Iu4W_pANXj9INd_0PvAGU1yivtGQqWFArEfHcMzgo_h4iV0U5dqu-kwfg",
    "content-length": body.length,
  },
};

const req = http.request(options, (res) => {
  console.log(`statusCode: ${res.statusCode}`);

  res.on('data', function (chunk) {
    console.log('Response: ' + chunk);
});
});

req.on("error", (error) => {
  console.error(error);
});

req.write(body)

req.end();
