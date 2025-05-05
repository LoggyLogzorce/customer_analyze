const reportBtn = document.getElementById('reportBtn');

document.getElementById("reportBtn").addEventListener('click', async (e) => {
    e.preventDefault();

    // Получаем базовые данные из DOM и localStorage
    const analysisData = JSON.parse(localStorage.getItem('lastAnalysis'));
    const filename = localStorage.getItem('filename');

    // Конвертируем графики в Base64
    const charts = {
        gender_pie: document.getElementById('genderChart')?.toDataURL('image/png'),
        age_histogram: document.getElementById('ageChart')?.toDataURL('image/png'),
    };

    // Формируем ReportData
    const reportData = {
        meta: {
            report_id: `report_${Date.now()}`,
            generated_at: new Date().toISOString(),
            filename: filename,
            total_users: analysisData?.demografi?.age_group?.Count || 0
        },

        demografi: {
            gender_distribution: analysisData?.demografi?.gender_dist || {},
            age_groups: analysisData?.demografi?.age_group || {},
            top_regions: analysisData?.demografi?.region_stat?.top5?.map(region => ({
                name: region.Name,
                count: region.Count
            })) || []
        },

        behavioral: {
            veterans: analysisData?.behavioral_analysis?.veterans || 0,
            newbies: analysisData?.behavioral_analysis?.newbies || 0,
            vips: {
                count: analysisData?.behavioral_analysis?.vips?.Count || 0,
                percentile: analysisData?.behavioral_analysis?.vips?.percentile75 || 0,
            }
        },

        // financials: {
        //     income: analysisData?.financials?.income || {},
        //     spending: analysisData?.financials?.spending || {},
        //     check_size_distribution: analysisData?.financials?.check_size_distribution || {}
        // },

        visualizations: {
            gender_pie: charts.gender_pie,
            age_histogram: charts.age_histogram,
        },

        // additional: {
        //     notes: document.getElementById('reportNotes')?.value || '',
        //     analyst_name: document.getElementById('analystName')?.value || 'Анонимный аналитик',
        //     segments: getSelectedSegments() // Ваша функция для получения сегментов
        // }
    };

    try {
        // Показываем лоадер
        showLoader();

        // Отправляем данные на сервер
        const response = await fetch('/api/v1/generate-report', {
            method: 'POST',
            headers: {'Content-Type': 'application/json'},
            body: JSON.stringify(reportData)
        });

        if (!response.ok) throw new Error('Ошибка генерации отчета');

        // Скачиваем файл
        const blob = await response.blob();
        downloadFile(blob, `report_${filename}.pdf`);

    } catch (error) {
        showError(`Ошибка: ${error.message}`);
    } finally {
        hideLoader();
    }
});

// Вспомогательные функции
function downloadFile(blob, filename) {
    const url = window.URL.createObjectURL(blob);
    const a = document.createElement('a');
    a.href = url;
    a.download = filename;
    document.body.appendChild(a);
    a.click();
    document.body.removeChild(a);
    URL.revokeObjectURL(url);
}

function showLoader() {
    reportBtn.classList.add('btn-loading');
    reportBtn.disabled = true;
}

function hideLoader() {
    reportBtn.classList.remove('btn-loading');
    reportBtn.disabled = false;
}

function showError(message) {
    const errorDiv = document.getElementById('reportError');
    errorDiv.textContent = message;
    errorDiv.style.display = 'block';
    setTimeout(() => errorDiv.style.display = 'none', 5000);
}