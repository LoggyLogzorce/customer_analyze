const resultsContainer = document.getElementById('results');
let genderChart = null;
let ageChart = null;
let map = null;

function displayResults(data) {
    resultsContainer.innerHTML = ''; // Очищаем предыдущие результаты

    // Основная информация
    const summary = document.createElement('div');
    summary.className = 'card mb-4';
    summary.innerHTML = `
        <div class="card-body" id='main_results'>
            <h4 class="card-title">📌 Основные результаты</h4>
            <p class="mb-1"><strong>Файл:</strong> ${data.filename}</p>
            <p class="mb-1"><strong>Статус:</strong> ${data.status}</p>
        </div>
    `;
    resultsContainer.appendChild(summary);

    // Демография
    if (data.demografi) {
        let card = []
        let content = []
        if (data.demografi.gender_dist) {
            content += [createSection('Распределение по полу', renderKeyValueGender(data.demografi.gender_dist))];
        }

        if (data.demografi.age_group) {
            content += [createSection('Возрастные группы', renderKeyValue(data.demografi.age_group))];
        }

        if (data.demografi.region_stat) {
            content += [createSection('Топ регионов', renderTopRegions(data.demografi.region_stat))];
        }
        card = createCard('📊 Демография', content);
        resultsContainer.appendChild(card);
    }

    // Поведенческий анализ
    if (data.behavioral_analysis) {
        let items = []
        if (data.behavioral_analysis.veterans) {
            items += [
                `<b>Ветераны:</b> ${data.behavioral_analysis.veterans} клиентов (регистрация до 2023 года)`,
            ];
        }
        if (data.behavioral_analysis.newbies) {
            if (items.length > 0) {
                items += [
                    `<br><b>Новички:</b> ${data.behavioral_analysis.newbies} клиентов (регистрация в 2025)`,
                ];
            } else {
                items += [
                    `<b>Новички:</b> ${data.behavioral_analysis.newbies} клиентов (регистрация в 2025)`,
                ];
            }
        }
        if (data.behavioral_analysis.vips) {
            if (items.length > 0) {
                items += [
                    `<br><b>VIP-клиенты:</b> ${data.behavioral_analysis.vips.Count} (перцентиль: ${data.behavioral_analysis.vips.percentile75.toFixed(2)})`
                ];
            } else {
                items += [
                    `<b>VIP-клиенты:</b> ${data.behavioral_analysis.vips.Count} (перцентиль: ${data.behavioral_analysis.vips.percentile75.toFixed(2)})`
                ];
            }
        }

        const card = createCard('📈 Поведенческий анализ', items);
        resultsContainer.appendChild(card);
    }

    // Очищаем предыдущие графики
    if (genderChart) genderChart.destroy();
    if (ageChart) ageChart.destroy();

    // Строим новые графики
    if (data.visualizations?.gender_pie) {
        renderGenderPie(data.visualizations.gender_pie);
    }

    if (data.visualizations?.age_group) {
        renderAgeHistogram(data.visualizations.age_group);
    }

    const reportBtn = document.getElementById('reportBtn')
    reportBtn.style.display = 'block'

    window.location.href = '#main_results'
}

// Вспомогательные функции
function createCard(title, content) {
    const card = document.createElement('div');
    card.className = 'card mb-4';
    card.innerHTML = `
        <div class="card-body">
            <h5 class="card-title">${title}</h5>
            ${Array.isArray(content) ? content.join('') : content}
        </div>
    `;
    return card;
}

function createSection(title, content) {
    return `
        <div class="mb-3">
            <h6 class="text-primary">${title}</h6>
            ${content}
        </div>
    `;
}

function renderKeyValueGender(data) {
    return `<ul class="list-unstyled">${
        Object.entries(data).map(([key, val]) => `
            <li><strong>${key}:</strong> ${val}</li>
        `).join('')
    }</ul>`;
}

function renderKeyValue(data) {
    return `<ul class="list-unstyled">${
        Object.entries(data).map(([key, val]) => `
            <li><strong>${key}:</strong> ${val} (${Math.round(val / data.Count * 100)}%)</li>
        `).join('')
    }</ul>`;
}

function renderTopRegions(data) {
    const top = data.top5 || data.top10;
    if (!top) return '';

    return `<ol class="list-group list-group-numbered">${
        top.map(region => `
            <li class="list-group-item d-flex justify-content-between align-items-start">
                <div class="ms-2 me-auto">${region.Name}</div>
                <span class="badge bg-primary rounded-pill">${region.Count}</span>
            </li>
        `).join('')
    }</ol>`;
}

// Круговая диаграмма по полу
function renderGenderPie(genderData) {
    const container = document.createElement('div');
    container.className = 'chart-container-cycle';

    const ctx = document.createElement('canvas');
    ctx.id = 'genderChart'
    container.appendChild(ctx);
    resultsContainer.appendChild(container);

    genderChart = new Chart(ctx, {
        type: 'doughnut',
        data: {
            labels: Object.keys(genderData),
            datasets: [{
                data: Object.values(genderData),
                backgroundColor: ['#36A2EB', '#FF6384']
            }]
        },
        options: {
            responsive: true,
            maintainAspectRatio: false,
            plugins: {
                title: {
                    display: true,
                    text: 'Распределение по полу',
                    font: { size: 14 }
                },
                legend: {
                    position: 'bottom',
                    labels: { font: { size: 12 } }
                }
            }
        }
    });
}

function renderAgeHistogram(ageData) {
    const container = document.createElement('div');
    container.className = 'chart-container-histogram';

    const ctx = document.createElement('canvas');
    ctx.id = 'ageChart'
    container.appendChild(ctx);
    resultsContainer.appendChild(container);

    ageChart = new Chart(ctx, {
        type: 'bar',
        data: {
            labels: Object.keys(ageData),
            datasets: [{
                label: 'Количество клиентов',
                data: Object.values(ageData),
                backgroundColor: '#4BC0C0'
            }]
        },
        options: {
            responsive: true,
            maintainAspectRatio: false,
            plugins: {
                title: {
                    display: true,
                    text: 'Возрастные группы',
                    font: { size: 14 }
                }
            },
            scales: {
                y: {
                    ticks: { font: { size: 12 } }
                },
                x: {
                    ticks: { font: { size: 12 } }
                }
            }
        }
    });
}