function loadQR() {
  fetch("/qrcode")
    .then(res => res.blob())
    .then(blob => {
      const url = URL.createObjectURL(blob);
      const img = document.getElementById("qr");
      img.src = url;
      img.style.display = "block";
    })
    .catch(err => {
      alert("Gagal ambil QR: " + err);
    });
}
