window.addEventListener("load", program);
function program() {
    document.getElementById("btn").addEventListener("click", async (e) => {
        e.preventDefault();

        const images = document.getElementById("input").files;
        const formData = createFormData(images);

        try {
            document.getElementById("idLoading").textContent = "Loading..."

            const response = await sendData(formData);
            const urlZip = await createUrlToDownloadImage(response);
            downloadImage(urlZip);

            document.getElementById("idLoading").textContent = ""
        } catch (err) {
            console.log(err);
            document.getElementById("idLoading").textContent = "Ha ocurrido un error."
        }
    });
}

function createFormData(images) {
    const formData = new FormData();

    for (const img of images) {
        formData.append("tif", img);
    }

    return formData;
}

async function sendData(formData) {
    return await fetch("http://localhost:8080/convert", {
        method: "POST",
        body: formData,
    });
}

async function createUrlToDownloadImage(response) {
    // blob? un formato raro que nos permitirá crear un objeto para generar un link
    return URL.createObjectURL(await response.blob());
}

function downloadImage(urlZip) {
    const a = document.createElement("a");
    a.href = urlZip;
    a.download = "Documentos.zip";
    document.body.appendChild(a);
    a.click();
    document.body.removeChild(a);
    URL.revokeObjectURL(urlZip);
}
