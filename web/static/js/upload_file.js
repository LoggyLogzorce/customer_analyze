// Обработка загрузки файла
async function handleUpload(event) {
    event.preventDefault();

    const form = event.target;
    const formData = new FormData(form);
    const progressBar = document.getElementById('progress-bar');
    const resultDiv = document.getElementById('upload-result');

    try {
        // Показать прогресс-бар
        progressBar.style.display = 'block';
        progressBar.innerHTML = 'Идет загрузка...';

        const response = await fetch('/api/v1/upload-file', {
            method: 'POST',
            body: formData
        });

        const result = await response.json();

        if (response.ok) {
            resultDiv.innerHTML = `
                        <div class="alert alert-success">
                            Файл успешно загружен!<br>
                            Имя файла: ${result.filename}<br>
                        </div>
                    `;
            localStorage.setItem("filename", result.filename)
            window.location.href = `/select-analysis/${encodeURIComponent(result.filename)}`
        } else {
            throw new Error(result.error || 'Ошибка загрузки');
        }
    } catch (error) {
        resultDiv.innerHTML = `
                    <div class="alert alert-danger">
                        Ошибка: ${error.message}
                    </div>
                `;
    } finally {
        progressBar.style.display = 'none';
        form.reset();
    }
}

// Показать информацию о выбранном файле
document.getElementById('fileInput').addEventListener('change', function(e) {
    const file = e.target.files[0];
    if (file) {
        document.getElementById('fileName').textContent = file.name;
        document.getElementById('fileSize').textContent =
            (file.size / 1024 / 1024).toFixed(2) + ' MB';
    }
});