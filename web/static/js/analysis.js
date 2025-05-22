document.addEventListener('DOMContentLoaded', () => {
    const form = document.getElementById('analysisForm');
    const submitBtn = document.getElementById('submitBtn');
    const resultsContainer = document.getElementById('results');

    // Обработчик отправки формы
    form.addEventListener('submit', async (e) => {
        e.preventDefault();

        // 1. Собираем данные формы
        const formData = new FormData(e.target);

        // 2. Добавляем имя файла
        const filename = localStorage.getItem('filename');
        formData.append('filename', filename);

        // 3. Удаляем отключенные элементы
        document.querySelectorAll('input:disabled, select:disabled').forEach(el => {
            if (el.name) formData.delete(el.name);
        });

        // 4. Логирование для отладки
        console.log('Отправляемые данные:');
        for (const [key, value] of formData.entries()) {
            console.log(key, '=', value);
        }

        // Показываем лоадер и блокируем кнопку
        submitBtn.classList.add('btn-loading');
        submitBtn.disabled = true;
        resultsContainer.innerHTML = '';

        try {
            // 5. Отправляем данные
            const response = await fetch('/api/v1/analyze', {
                method: 'POST',
                body: formData // Отправляем FormData напрямую
            });

            const data = await response.json()
            if (!response.ok) {
                throw new Error(data.error);
            }
            localStorage.setItem('lastAnalysis', JSON.stringify(data))
            displayResults(data);

        } catch (error) {
            // Показываем ошибку
            resultsContainer.innerHTML = `
                <div class="alert alert-danger">
                    Ошибка: ${error.message}
                </div>
            `;
            console.error('Ошибка:', error);
        } finally {
            // Скрываем лоадер и разблокируем кнопку
            submitBtn.classList.remove('btn-loading');
            submitBtn.disabled = false;
        }
    });
});