package templates

const HTMLTemplate = `
<!DOCTYPE html>
<html lang="ru">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Load Tester</title>
    <style>
        :root {
            /* Светлая тема - нежно-синие тона */
            --bg-primary: #f0f4ff;
            --bg-secondary: #ffffff;
            --bg-card: #fafbff;
            --text-primary: #1e3a8a;
            --text-secondary: #3b82f6;
            --text-muted: #64748b;
            --border-color: #e2e8f0;
            --accent-color: #3b82f6;
            --accent-hover: #2563eb;
            --success-color: #10b981;
            --success-bg: #ecfdf5;
            --error-color: #ef4444;
            --error-bg: #fef2f2;
            --warning-color: #f59e0b;
            --warning-bg: #fffbeb;
            --shadow: 0 4px 6px -1px rgba(0, 0, 0, 0.1);
            --shadow-lg: 0 10px 15px -3px rgba(0, 0, 0, 0.1);
        }

        [data-theme="dark"] {
            /* Темная тема - темно-синие тона */
            --bg-primary: #0f172a;
            --bg-secondary: #1e293b;
            --bg-card: #334155;
            --text-primary: #e2e8f0;
            --text-secondary: #94a3b8;
            --text-muted: #64748b;
            --border-color: #475569;
            --accent-color: #3b82f6;
            --accent-hover: #60a5fa;
            --success-color: #10b981;
            --success-bg: #064e3b;
            --error-color: #ef4444;
            --error-bg: #7f1d1d;
            --warning-color: #f59e0b;
            --warning-bg: #78350f;
            --shadow: 0 4px 6px -1px rgba(0, 0, 0, 0.3);
            --shadow-lg: 0 10px 15px -3px rgba(0, 0, 0, 0.3);
        }

        * {
            margin: 0;
           
            box-sizing: border-box;
        }

        body {
            font-family: 'Segoe UI', Tahoma, Geneva, Verdana, sans-serif;
            background: var(--bg-primary);
            color: var(--text-primary);
            transition: all 0.3s ease;
            min-height: 100vh;
        }

        .theme-toggle {
            background: var(--bg-card);
            border: 1px solid var(--border-color);
            border-radius: 20px;
            padding: 6px 12px;
            cursor: pointer;
            display: flex;
            align-items: center;
            gap: 6px;
            font-size: 13px;
            color: var(--text-primary);
            transition: all 0.3s ease;
        }

        .theme-toggle:hover {
            background: var(--accent-color);
            color: white;
            border-color: var(--accent-color);
        }

        .theme-icon {
            font-size: 14px;
        }

        .navigation {
            background: var(--bg-secondary);
            border-bottom: 1px solid var(--border-color);
            box-shadow: var(--shadow);
            position: sticky;
            top: 0;
            z-index: 100;
        }

        .nav-container {
            max-width: 1200px;
            margin: 0 auto;
            padding: 0 20px;
            display: flex;
            justify-content: space-between;
            align-items: center;
            height: 60px;
        }

        .nav-left {
            display: flex;
            align-items: center;
            gap: 40px;
        }

        .nav-brand {
            font-size: 20px;
            font-weight: 700;
            color: var(--accent-color);
        }

        .nav-menu {
            display: flex;
        }

        .nav-right {
            display: flex;
            align-items: center;
        }

        @media (max-width: 768px) {
            .nav-container {
                flex-direction: column;
                height: auto;
                padding: 16px 20px;
                gap: 16px;
            }
            
            .nav-left {
                flex-direction: column;
                gap: 16px;
                width: 100%;
            }
            
            .nav-menu {
                gap: 20px;
                justify-content: center;
            }
            
            .nav-right {
                width: 100%;
                justify-content: center;
            }
        }

        .nav-link {
            color: var(--text-muted);
            text-decoration: none;
            font-weight: 500;
            padding: 8px 16px;
            border-radius: 6px;
            transition: all 0.3s ease;
			margin: 3px;
        }

        .nav-link:hover {
            color: var(--accent-color);
            background: var(--bg-card);
        }

        .nav-link.active {
            color: var(--accent-color);
            background: var(--bg-card);
        }

        .agreement-content {
            max-width: 800px;
            margin: 0 auto;
            line-height: 1.7;
        }

        .instructions-content {
            max-width: 900px;
            margin: 0 auto;
            line-height: 1.7;
        }

        .intro-section, .settings-section, .network-section, .examples-section, .results-section, .tips-section {
            margin-bottom: 40px;
        }

        .intro-section h2, .settings-section h2, .network-section h2, .examples-section h2, .results-section h2, .tips-section h2 {
            color: var(--accent-color);
            margin-bottom: 20px;
            font-size: 24px;
            border-bottom: 2px solid var(--accent-color);
            padding-bottom: 8px;
        }

        .setting-item, .result-item, .tip-item {
            background: var(--bg-card);
            border: 1px solid var(--border-color);
            border-radius: 12px;
            padding: 24px;
            margin-bottom: 20px;
        }

        .setting-item h3, .result-item h3, .tip-item h3 {
            color: var(--text-primary);
            margin-bottom: 12px;
            font-size: 18px;
        }

        .setting-item p, .result-item p, .tip-item p {
            color: var(--text-primary);
            margin-bottom: 12px;
        }

        .setting-item ul, .result-item ul, .tip-item ul {
            margin: 12px 0;
            padding-left: 20px;
        }

        .setting-item li, .result-item li, .tip-item li {
            color: var(--text-primary);
            margin-bottom: 6px;
        }

        .setting-item code {
            background: var(--bg-primary);
            padding: 2px 6px;
            border-radius: 4px;
            font-family: 'Courier New', monospace;
            color: var(--accent-color);
        }

        .performance-table, .network-table {
            background: var(--success-bg);
            border: 1px solid var(--success-color);
            border-radius: 8px;
            padding: 16px;
            margin: 16px 0;
            overflow-wrap: break-word;
            word-wrap: break-word;
        }

        .performance-table h4, .network-table h3 {
            color: var(--success-color);
            margin-bottom: 12px;
            line-height: 1.4;
            word-wrap: break-word;
        }

        .info-block {
            background: var(--bg-secondary);
            border: 1px solid var(--border-color);
            border-radius: 8px;
            padding: 16px;
            margin-top: 16px;
        }

        .info-block h4 {
            color: var(--accent-color);
            margin-bottom: 8px;
        }

        .config-example {
            margin-bottom: 24px;
        }

        .config-example h3 {
            margin-bottom: 12px;
        }

        .config-box {
            background: var(--bg-secondary);
            border: 1px solid var(--border-color);
            border-radius: 8px;
            padding: 16px;
            margin: 12px 0;
            font-family: 'Courier New', monospace;
        }

        .config-box p {
            margin: 4px 0;
            color: var(--text-primary);
        }

        .warning-small {
            background: var(--warning-bg);
            color: var(--warning-color);
            padding: 8px 12px;
            border-radius: 6px;
            font-size: 14px;
            margin-top: 12px;
        }

        .warning-block {
            background: var(--error-bg);
            border: 2px solid var(--error-color);
            border-radius: 12px;
            padding: 24px;
            margin-bottom: 32px;
        }

        .warning-block h2 {
            color: var(--error-color);
            margin-bottom: 16px;
            font-size: 20px;
        }

        .warning-block p {
            color: var(--text-primary);
            font-size: 16px;
        }

        .info-section, .purpose-section, .legal-section, .license-section {
            margin-bottom: 32px;
        }

        .info-section h2, .purpose-section h2, .legal-section h2, .license-section h2 {
            color: var(--accent-color);
            margin-bottom: 16px;
            font-size: 20px;
        }

        .info-section h3, .purpose-section h3, .legal-section h3 {
            color: var(--text-primary);
            margin: 20px 0 12px 0;
            font-size: 16px;
        }

        .info-section p, .purpose-section p, .legal-section p, .license-section p {
            color: var(--text-primary);
            margin-bottom: 16px;
        }

        .info-section ul, .purpose-section ul, .legal-section ul {
            margin: 16px 0;
            padding-left: 24px;
        }

        .info-section li, .purpose-section li, .legal-section li {
            color: var(--text-primary);
            margin-bottom: 8px;
        }

        .responsibility-block {
            background: var(--warning-bg);
            border: 1px solid var(--warning-color);
            border-radius: 8px;
            padding: 20px;
            margin: 20px 0;
        }

        .responsibility-block h3 {
            color: var(--warning-color);
            margin-bottom: 12px;
        }

        .responsibility-block ul {
            margin: 12px 0;
            padding-left: 20px;
        }

        .responsibility-block li {
            color: var(--text-primary);
            margin-bottom: 6px;
        }

        .agreement-footer {
            background: var(--bg-card);
            border: 1px solid var(--border-color);
            border-radius: 12px;
            padding: 24px;
            margin-top: 32px;
        }

        .agreement-footer p {
            font-weight: 600;
            margin-bottom: 16px;
            color: var(--text-primary);
        }

        .agreement-footer ul {
            margin: 0;
            padding-left: 20px;
        }

        .agreement-footer li {
            color: var(--text-primary);
            margin-bottom: 8px;
        }

        .cookie-notice {
            position: fixed;
            bottom: 20px;
            left: 20px;
            right: 20px;
            background: var(--bg-secondary);
            border: 1px solid var(--border-color);
            border-radius: 12px;
            padding: 16px;
            box-shadow: var(--shadow-lg);
            z-index: 1000;
            transform: translateY(100px);
            opacity: 0;
            transition: all 0.3s ease;
        }

        .cookie-notice.show {
            transform: translateY(0);
            opacity: 1;
        }

        .cookie-notice-content {
            display: flex;
            justify-content: space-between;
            align-items: center;
            gap: 16px;
        }

        .cookie-notice-text {
            font-size: 14px;
            color: var(--text-muted);
        }

        .cookie-notice-button {
            background: var(--accent-color);
            color: white;
            border: none;
            padding: 8px 16px;
            border-radius: 6px;
            cursor: pointer;
            font-size: 14px;
            transition: background-color 0.3s ease;
        }

        .cookie-notice-button:hover {
            background: var(--accent-hover);
        }

        .container {
            max-width: 1200px;
            margin: 0 auto;
            padding: 20px 20px 40px 20px;
        }

        .main-card {
            background: var(--bg-secondary);
            border-radius: 16px;
            padding: 40px;
            box-shadow: var(--shadow-lg);
            border: 1px solid var(--border-color);
        }

        h1 {
            text-align: center;
            margin-bottom: 40px;
            font-size: 2.5rem;
            background: linear-gradient(135deg, var(--accent-color), #60a5fa);
            -webkit-background-clip: text;
            background-clip: text;
        }

        .form-group {
            margin-bottom: 24px;
        }

        label {
            display: block;
            margin-bottom: 8px;
            font-weight: 600;
            color: var(--text-primary);
            font-size: 14px;
        }

        .field-description {
            font-size: 12px;
            color: var(--text-muted);
            margin-top: 6px;
            line-height: 1.5;
        }

        input, select {
            width: 100%;
            padding: 12px 16px;
            border: 2px solid var(--border-color);
            border-radius: 8px;
            font-size: 16px;
            background: var(--bg-card);
            color: var(--text-primary);
            transition: all 0.3s ease;
        }

        input:focus, select:focus {
            border-color: var(--accent-color);
            outline: none;
            box-shadow: 0 0 0 3px rgba(59, 130, 246, 0.1);
        }

        .form-row {
            display: grid;
            grid-template-columns: 1fr 1fr;
            gap: 24px;
        }

        @media (max-width: 768px) {
            .form-row {
                grid-template-columns: 1fr;
            }
            
            .checkbox-container {
                flex-direction: column;
                gap: 8px;
                text-align: center;
                padding: 16px;
                max-width: none;
                margin: 0 16px;
            }
            
            .buttons {
                flex-direction: column;
                align-items: center;
                gap: 12px;
            }
            
            button {
                width: 100%;
                max-width: 300px;
            }
            
            .instructions-content {
                max-width: 100%;
                padding: 0 16px;
            }
            
            .setting-item, .result-item, .tip-item {
                padding: 16px;
                margin-bottom: 16px;
            }
            
            .network-table h3, .performance-table h4 {
                font-size: 16px;
                line-height: 1.4;
                word-wrap: break-word;
            }
            
            .config-box {
                padding: 12px;
                font-size: 14px;
            }
            
            .intro-section h2, .settings-section h2, .network-section h2, 
            .examples-section h2, .results-section h2, .tips-section h2 {
                font-size: 20px;
                line-height: 1.3;
            }
        }

        .agreement-checkbox {
            text-align: center;
            margin: 32px 0 24px 0;
        }

        .checkbox-container {
            display: inline-flex;
            align-items: center;
            gap: 12px;
            cursor: pointer;
            font-size: 14px;
            color: var(--text-primary);
            padding: 16px 24px;
            background: var(--bg-card);
            border: 1px solid var(--border-color);
            border-radius: 12px;
            transition: all 0.3s ease;
            max-width: 500px;
            text-align: left;
        }

        .checkbox-container:hover {
            border-color: var(--accent-color);
            box-shadow: 0 0 0 3px rgba(59, 130, 246, 0.1);
        }

        .checkbox-container input[type="checkbox"] {
            display: none;
        }

        .checkmark {
            width: 20px;
            height: 20px;
            border: 2px solid var(--border-color);
            border-radius: 4px;
            background: var(--bg-secondary);
            position: relative;
            transition: all 0.3s ease;
            flex-shrink: 0;
        }

        .checkbox-container input[type="checkbox"]:checked + .checkmark {
            background: var(--accent-color);
            border-color: var(--accent-color);
        }

        .checkbox-container input[type="checkbox"]:checked + .checkmark::after {
            content: '✓';
            position: absolute;
            top: 50%;
            left: 50%;
            transform: translate(-50%, -50%);
            color: white;
            font-size: 14px;
            font-weight: bold;
        }

        .checkbox-text {
            line-height: 1.5;
        }

        .agreement-link {
            color: var(--accent-color);
            text-decoration: none;
            font-weight: 600;
            transition: color 0.3s ease;
        }

        .agreement-link:hover {
            color: var(--accent-hover);
            text-decoration: underline;
        }

        .buttons {
            text-align: center;
            margin: 24px 0 40px 0;
            display: flex;
            gap: 16px;
            justify-content: center;
        }

        button {
            padding: 14px 32px;
            border: none;
            border-radius: 8px;
            font-size: 16px;
            font-weight: 600;
            cursor: pointer;
            transition: all 0.3s ease;
            display: flex;
            align-items: center;
            gap: 8px;
        }

        .start-btn {
            background: linear-gradient(135deg, var(--success-color), #34d399);
            color: white;
        }

        .start-btn:hover:not(:disabled) {
            transform: translateY(-2px);
            box-shadow: var(--shadow-lg);
        }

        .stop-btn {
            background: linear-gradient(135deg, var(--error-color), #f87171);
            color: white;
        }

        .stop-btn:hover:not(:disabled) {
            transform: translateY(-2px);
            box-shadow: var(--shadow-lg);
        }

        button:disabled {
            background: var(--text-muted);
            cursor: not-allowed;
            transform: none;
        }

        .status {
            text-align: center;
            padding: 20px;
            margin: 24px 0;
            border-radius: 12px;
            font-weight: 600;
            font-size: 18px;
            transition: all 0.3s ease;
        }

        .status.ready {
            background: var(--success-bg);
            color: var(--success-color);
            border: 1px solid var(--success-color);
        }

        .status.running {
            background: var(--warning-bg);
            color: var(--warning-color);
            border: 1px solid var(--warning-color);
            animation: pulse 2s infinite;
        }

        .status.completed {
            background: var(--bg-card);
            color: var(--text-primary);
            border: 1px solid var(--border-color);
        }

        @keyframes pulse {
            0%, 100% { opacity: 1; }
            50% { opacity: 0.8; }
        }

        .progress-container {
            margin: 24px 0;
        }

        .progress-bar {
            width: 100%;
            height: 24px;
            background: var(--bg-card);
            border-radius: 12px;
            overflow: hidden;
            border: 1px solid var(--border-color);
        }

        .progress-fill {
            height: 100%;
            background: linear-gradient(90deg, var(--accent-color), #60a5fa);
            transition: width 0.3s ease;
            border-radius: 12px;
            position: relative;
        }

        .progress-fill::after {
            content: '';
            position: absolute;
            top: 0;
            left: 0;
            right: 0;
            bottom: 0;
            background: linear-gradient(90deg, transparent, rgba(255,255,255,0.3), transparent);
            animation: shimmer 2s infinite;
        }

        @keyframes shimmer {
            0% { transform: translateX(-100%); }
            100% { transform: translateX(100%); }
        }

        .stats {
            display: grid;
            grid-template-columns: repeat(auto-fit, minmax(175px, 1fr));
            gap: 20px;
            margin-top: 40px;
        }

        .stat-card {
            background: var(--bg-card);
            padding: 24px;
            border-radius: 12px;
            border: 1px solid var(--border-color);
            transition: all 0.3s ease;
            position: relative;
            overflow: hidden;
        }

        .stat-card::before {
            content: '';
            position: absolute;
            top: 0;
            left: 0;
            right: 0;
            height: 4px;
            background: var(--accent-color);
        }

        .stat-card:hover {
            transform: translateY(-4px);
            box-shadow: var(--shadow-lg);
        }

        .stat-card.error::before {
            background: var(--error-color);
        }

        .stat-card.warning::before {
            background: var(--warning-color);
        }

        .stat-card.success::before {
            background: var(--success-color);
        }

        .stat-title {
            font-size: 14px;
            color: var(--text-muted);
            margin-bottom: 8px;
            font-weight: 500;
        }

        .stat-value {
            font-size: 28px;
            font-weight: 700;
            color: var(--text-primary);
            margin-bottom: 4px;
        }

        .stat-subtitle {
            font-size: 12px;
            color: var(--text-muted);
        }
    </style>
</head>
<body>
    <div class="cookie-notice" id="cookieNotice">
        <div class="cookie-notice-content">
            <div class="cookie-notice-text">
                Мы используем localStorage для сохранения ваших настроек темы. Продолжая использование, вы соглашаетесь с этим.
            </div>
            <button class="cookie-notice-button" onclick="acceptCookies()">Понятно</button>
        </div>
    </div>

    <nav class="navigation">
        <div class="nav-container">
            <div class="nav-left">
                <div class="nav-brand">🚀 Load Tester</div>
                <div class="nav-menu">
                    <a href="#" class="nav-link active" onclick="showPage('main')">Главная</a>
                    <a href="#" class="nav-link" onclick="showPage('instructions')">Инструкция</a>
                    <a href="#" class="nav-link" onclick="showPage('agreement')">Соглашение</a>
                </div>
            </div>
            <div class="nav-right">
                <div class="theme-toggle" onclick="toggleTheme()">
                    <span class="theme-icon" id="themeIcon">🌙</span>
                    <span id="themeText">Темная</span>
                </div>
            </div>
        </div>
    </nav>

    <div class="container">
        <div class="main-card" id="mainPage">
            <h1>⚙️ Настройки тестирования</h1>
            
            <form id="testForm">
                <div class="form-group">
                    <label for="targetURL">URL цели:</label>
                    <input type="url" id="targetURL" name="targetURL" placeholder="https://site.ru" required>
                    <div class="field-description">Полный URL сервера для нагрузочного тестирования (включая протокол http:// или https://)</div>
                </div>
                
                <div class="form-row">
                    <div class="form-group">
                        <label for="method">HTTP метод:</label>
                        <select id="method" name="method">
                            <option value="GET">GET</option>
                            <option value="POST">POST</option>
                            <option value="PUT">PUT</option>
                            <option value="DELETE">DELETE</option>
                            <option value="HEAD">HEAD</option>
                            <option value="OPTIONS">OPTIONS</option>
                        </select>
                        <div class="field-description">Тип HTTP запроса. GET - для чтения данных, POST - для отправки данных</div>
                    </div>
                    <div class="form-group">
                        <label for="totalRequests">Всего запросов:</label>
                        <input type="number" id="totalRequests" name="totalRequests" value="300000" min="1" required>
                        <div class="field-description">Общее количество HTTP запросов, которые будут выполнены во время теста</div>
                    </div>
                </div>
                
                <div class="form-row">
                    <div class="form-group">
                        <label for="maxConcurrency">Конкурентность:</label>
                        <input type="number" id="maxConcurrency" name="maxConcurrency" value="10000" min="1" required>
                        <div class="field-description">Максимальное количество одновременных запросов. Больше = выше нагрузка</div>
                    </div>
                    <div class="form-group">
                        <label for="timeout">Таймаут (сек):</label>
                        <input type="number" id="timeout" name="timeout" value="2" min="1" required>
                        <div class="field-description">Время ожидания ответа от сервера. Запрос отменяется при превышении</div>
                    </div>
                </div>
                
                <div class="form-group">
                    <label for="maxMemory">Максимальная память (GB):</label>
                    <input type="number" id="maxMemory" name="maxMemory" value="30" min="1" step="0.1" required>
                    <div class="field-description">Лимит памяти для программы в гигабайтах. Рекомендуется 70-80% от доступной RAM на вашем ПК или хосте, где запускаете программу<br>Если вы выделите программе больше памяти чем у вас есть, вы можете столкнуться с зависанием</div>
                </div>
            </form>
            
            <div class="agreement-checkbox">
                <label class="checkbox-container">
                    <input type="checkbox" id="agreementCheckbox" onchange="toggleStartButton()">
                    <span class="checkmark"></span>
                    <span class="checkbox-text">
                        Принимаю <a href="#" onclick="showPage('agreement')" class="agreement-link">пользовательское соглашение</a> 
                        и понимаю принцип работы программы
                    </span>
                </label>
            </div>
            
            <div class="buttons">
                <button type="button" class="start-btn" id="startBtn" onclick="startTest()" disabled>
                    <span>🚀</span> Запустить тест
                </button>
                <button type="button" class="stop-btn" id="stopBtn" onclick="stopTest()" disabled>
                    <span>⏹️</span> Остановить
                </button>
            </div>
            
            <div id="status" class="status ready">Готов к запуску</div>
            
            <div class="progress-container">
                <div class="progress-bar">
                    <div class="progress-fill" id="progressFill" style="width: 0%"></div>
                </div>
            </div>
            
            <div class="stats" id="stats" style="display: none;">
                <div class="stat-card success">
                    <div class="stat-title">Всего запросов</div>
					<div class="stat-subtitle">Сколько запросов уже сделал тест после начала запуска</div>
                    <br>
					<div class="stat-value" id="totalStat">0</div>
                    <div class="stat-subtitle">RPS: <span id="rpsStat">0</span></div>
                </div>
                <div class="stat-card success">
                    <div class="stat-title">Успешных (200 OK)</div>
					<div class="stat-subtitle">Сколько запросов сервер успел обработать перед тем как завис</div>
                    <br>
					<div class="stat-value" id="successStat">0</div>
                    <div class="stat-subtitle"><span id="successPercent">0</span>%</div>
                </div>
                <div class="stat-card error">
                    <div class="stat-title">HTTP ошибки</div>
					<div class="stat-subtitle">Ошибки которыми начал отвечать сервер</div>
                    <br>
					<div class="stat-value" id="failsStat">0</div>
                    <div class="stat-subtitle">4xx: <span id="status4xxStat">0</span> | 5xx: <span id="status5xxStat">0</span></div>
                </div>
                <div class="stat-card warning">
                    <div class="stat-title">Сетевые ошибки</div>
					<div class="stat-subtitle">Это уже перегрузка канала нагрузочного хоста</div>
                    <br>
					<div class="stat-value" id="errorsStat">0</div>
                    <div class="stat-subtitle">Таймауты: <span id="timeoutsStat">0</span> | Соединения: <span id="connectionsStat">0</span></div>
                </div>
            </div>
        </div>

        <div class="main-card" id="instructionsPage" style="display: none;">
            <h1>📖 Инструкция по использованию</h1>
            
            <div class="instructions-content">
                <div class="intro-section">
                    <h2>🎯 Что такое Load Tester?</h2>
                    <p>Load Tester — это профессиональный инструмент для нагрузочного тестирования веб-серверов. Программа создает высокую нагрузку на HTTP-сервер путем запуска большого количества одновременных Go-рутин (горутин), которые выполняют HTTP-запросы к указанному URL.</p>
                </div>

                <div class="settings-section">
                    <h2>⚙️ Описание настроек</h2>
                    
                    <div class="setting-item">
                        <h3>🌐 URL цели</h3>
                        <p><strong>Что это:</strong> Полный адрес сервера для тестирования</p>
                        <p><strong>Примеры:</strong></p>
                        <ul>
                            <li><code>https://example.com</code> — тестирование главной страницы</li>
                            <li><code>https://api.example.com/health</code> — тестирование API endpoint</li>
                            <li><code>http://localhost:3000</code> — тестирование локального сервера</li>
                        </ul>
                        <p><strong>Важно:</strong> Обязательно указывайте протокол (http:// или https://)</p>
                    </div>

                    <div class="setting-item">
                        <h3>📡 HTTP метод</h3>
                        <p><strong>GET</strong> — для чтения данных (по умолчанию)</p>
                        <p><strong>POST</strong> — для отправки данных на сервер</p>
                        <p><strong>PUT/DELETE</strong> — для REST API тестирования</p>
                        <p><strong>HEAD</strong> — только заголовки, без тела ответа</p>
                        <p><strong>OPTIONS</strong> — для проверки поддерживаемых методов</p>
                    </div>

                    <div class="setting-item">
                        <h3>🔢 Всего запросов</h3>
                        <p><strong>Что это:</strong> Общее количество HTTP-запросов для выполнения</p>
                        <p><strong>Рекомендации:</strong></p>
                        <ul>
                            <li><strong>Легкий тест:</strong> 1,000 - 10,000 запросов</li>
                            <li><strong>Средний тест:</strong> 50,000 - 100,000 запросов</li>
                            <li><strong>Интенсивный тест:</strong> 300,000+ запросов</li>
                        </ul>
                    </div>

                    <div class="setting-item">
                        <h3>⚡ Конкурентность</h3>
                        <p><strong>Что это:</strong> Количество одновременных соединений</p>
                        <p><strong>Влияние на нагрузку:</strong> Чем больше — тем выше нагрузка</p>
                        <div class="performance-table">
                            <h4>Рекомендации по процессорам:</h4>
                            <ul>
                                <li><strong>2-4 ядра (бюджетные ПК):</strong> 100-1,000 горутин</li>
                                <li><strong>4-8 ядер (средние ПК):</strong> 1,000-5,000 горутин</li>
                                <li><strong>8-16 ядер (мощные ПК):</strong> 5,000-15,000 горутин</li>
                                <li><strong>16+ ядер (серверы):</strong> 15,000-50,000+ горутин</li>
                            </ul>
                        </div>
                    </div>

                    <div class="setting-item">
                        <h3>⏱️ Таймаут</h3>
                        <p><strong>Что это:</strong> Время ожидания ответа от сервера</p>
                        <p><strong>Рекомендации:</strong></p>
                        <ul>
                            <li><strong>1-2 сек:</strong> Для быстрых API и статических страниц</li>
                            <li><strong>3-5 сек:</strong> Для обычных веб-приложений</li>
                            <li><strong>10+ сек:</strong> Для медленных или перегруженных серверов</li>
                        </ul>
                    </div>

                    <div class="setting-item">
                        <h3>💾 Максимальная память</h3>
                        <p><strong>Что это:</strong> Лимит оперативной памяти для программы</p>
                        <p><strong>Рекомендации:</strong></p>
                        <ul>
                            <li><strong>8 GB RAM:</strong> Установите 4-6 GB</li>
                            <li><strong>16 GB RAM:</strong> Установите 10-12 GB</li>
                            <li><strong>32 GB RAM:</strong> Установите 20-25 GB</li>
                            <li><strong>64+ GB RAM:</strong> Установите 40-50 GB</li>
                        </ul>
                        <div class="warning-small">
                            ⚠️ Не устанавливайте больше 80% от доступной RAM
                        </div>
                    </div>
                </div>

                <div class="network-section">
                    <h2>🌐 Интернет-канал</h2>
                    
                    <div class="network-table">
                        <h3>Рекомендации по скорости интернета:</h3>
                        <ul>
                            <li><strong>100 Мбит/с:</strong> До 1,000 одновременных соединений</li>
                            <li><strong>500 Мбит/с:</strong> До 5,000 одновременных соединений</li>
                            <li><strong>1 Гбит/с:</strong> До 10,000 одновременных соединений</li>
                            <li><strong>10 Гбит/с:</strong> До 50,000+ одновременных соединений</li>
                        </ul>
                        
                        <div class="info-block">
                            <h4>💡 Важные факторы:</h4>
                            <ul>
                                <li><strong>Пинг:</strong> Чем меньше, тем лучше (оптимально < 50ms)</li>
                                <li><strong>Стабильность:</strong> Избегайте Wi-Fi для серьезных тестов</li>
                                <li><strong>Провайдер:</strong> Некоторые блокируют массовые запросы</li>
                            </ul>
                        </div>
                    </div>
                </div>

                <div class="examples-section">
                    <h2>📋 Примеры конфигураций</h2>
                    
                    <div class="config-example">
                        <h3>🟢 Легкий тест (проверка доступности)</h3>
                        <div class="config-box">
                            <p><strong>URL:</strong> https://your-site.com</p>
                            <p><strong>Метод:</strong> GET</p>
                            <p><strong>Запросов:</strong> 1,000</p>
                            <p><strong>Конкурентность:</strong> 100</p>
                            <p><strong>Таймаут:</strong> 5 сек</p>
                            <p><strong>Память:</strong> 2 GB</p>
                        </div>
                        <p><em>Подходит для проверки базовой работоспособности сайта</em></p>
                    </div>

                    <div class="config-example">
                        <h3>🟡 Средний тест (реальная нагрузка)</h3>
                        <div class="config-box">
                            <p><strong>URL:</strong> https://api.your-site.com/users</p>
                            <p><strong>Метод:</strong> GET</p>
                            <p><strong>Запросов:</strong> 50,000</p>
                            <p><strong>Конкурентность:</strong> 2,000</p>
                            <p><strong>Таймаут:</strong> 3 сек</p>
                            <p><strong>Память:</strong> 8 GB</p>
                        </div>
                        <p><em>Имитирует реальную нагрузку от пользователей</em></p>
                    </div>

                    <div class="config-example">
                        <h3>🔴 Стресс-тест (максимальная нагрузка)</h3>
                        <div class="config-box">
                            <p><strong>URL:</strong> https://your-site.com</p>
                            <p><strong>Метод:</strong> GET</p>
                            <p><strong>Запросов:</strong> 500,000</p>
                            <p><strong>Конкурентность:</strong> 10,000</p>
                            <p><strong>Таймаут:</strong> 2 сек</p>
                            <p><strong>Память:</strong> 20 GB</p>
                        </div>
                        <p><em>Проверяет пределы производительности сервера</em></p>
                    </div>
                </div>

                <div class="results-section">
                    <h2>📊 Интерпретация результатов</h2>
                    
                    <div class="result-item">
                        <h3>✅ Успешные запросы (200 OK)</h3>
                        <p>Количество запросов, которые сервер успешно обработал</p>
                        <ul>
                            <li><strong>95-100%:</strong> Отличная производительность</li>
                            <li><strong>80-95%:</strong> Хорошая производительность</li>
                            <li><strong>< 80%:</strong> Есть проблемы с производительностью</li>
                        </ul>
                    </div>

                    <div class="result-item">
                        <h3>❌ HTTP ошибки (4xx/5xx)</h3>
                        <p><strong>4xx ошибки:</strong> Проблемы с запросами (неверный URL, авторизация)</p>
                        <p><strong>5xx ошибки:</strong> Проблемы сервера (перегрузка, сбои)</p>
                    </div>

                    <div class="result-item">
                        <h3>🌐 Сетевые ошибки</h3>
                        <p><strong>Таймауты:</strong> Сервер не успевает отвечать</p>
                        <p><strong>Соединения:</strong> Проблемы с подключением к серверу</p>
                    </div>

                    <div class="result-item">
                        <h3>⚡ RPS (Requests Per Second)</h3>
                        <p>Количество запросов в секунду — показатель производительности</p>
                        <ul>
                            <li><strong>1,000+ RPS:</strong> Высокая производительность</li>
                            <li><strong>100-1,000 RPS:</strong> Средняя производительность</li>
                            <li><strong>< 100 RPS:</strong> Низкая производительность</li>
                        </ul>
                    </div>
                </div>

                <div class="tips-section">
                    <h2>💡 Советы и рекомендации</h2>
                    
                    <div class="tip-item">
                        <h3>🎯 Перед тестированием</h3>
                        <ul>
                            <li>Получите разрешение владельца сервера</li>
                            <li>Начните с малых нагрузок</li>
                            <li>Проверьте стабильность интернет-соединения</li>
                            <li>Закройте ненужные программы для освобождения ресурсов</li>
                        </ul>
                    </div>

                    <div class="tip-item">
                        <h3>⚡ Оптимизация производительности</h3>
                        <ul>
                            <li>Используйте проводное подключение вместо Wi-Fi</li>
                            <li>Запускайте тесты в нерабочее время</li>
                            <li>Мониторьте использование CPU и RAM</li>
                            <li>Постепенно увеличивайте нагрузку</li>
                        </ul>
                    </div>

                    <div class="tip-item">
                        <h3>🛡️ Безопасность</h3>
                        <ul>
                            <li>Не тестируйте чужие сайты без разрешения</li>
                            <li>Соблюдайте законодательство вашей страны</li>
                            <li>Используйте только для легальных целей</li>
                            <li>Будьте готовы остановить тест при необходимости</li>
                        </ul>
                    </div>
                </div>
            </div>
        </div>

        <div class="main-card" id="agreementPage" style="display: none;">
            <h1>📋 Пользовательское соглашение</h1>
            
            <div class="agreement-content">
                <div class="warning-block">
                    <h2>⚠️ Важное предупреждение</h2>
                    <p><strong>Вы используете программу на свой страх и риск.</strong> Автор программы не несет ответственности за всякий причиненный вред, который может наступить в процессе использования данного программного обеспечения.</p>
                </div>

                <div class="info-section">
                    <h2>🔧 Принцип работы программы</h2>
                    <p>Программа создает сильную нагрузку на HTTP-сервер путем запуска большого количества одновременных потоков Go Routine (горутин), которые выполняют попытки соединения и получения ответа на указанный URL-адрес.</p>
                    
                    <h3>Факторы, влияющие на силу нагрузки:</h3>
                    <ul>
                        <li><strong>Количество параллельных потоков</strong> — зависит от мощности вашего ПК, в основном от процессора</li>
                        <li><strong>Пропускная способность интернет-соединения</strong> — влияет на скорость отправки запросов</li>
                        <li><strong>Настройки конкурентности</strong> — количество одновременных соединений</li>
                        <li><strong>Таймауты и лимиты памяти</strong> — определяют стабильность работы</li>
                    </ul>
                </div>

                <div class="purpose-section">
                    <h2>🎯 Назначение программы</h2>
                    <p>Программа предназначена для тестирования веб-серверов и веб-инфраструктуры на отказоустойчивость с целью последующего выявления слабых мест сервиса и их устранения.</p>
                    
                    <h3>Легальное использование:</h3>
                    <ul>
                        <li>Тестирование собственных серверов и приложений</li>
                        <li>Нагрузочное тестирование с разрешения владельца ресурса</li>
                        <li>Образовательные цели и изучение производительности</li>
                        <li>Оптимизация инфраструктуры и выявление узких мест</li>
                    </ul>
                </div>

                <div class="legal-section">
                    <h2>⚖️ Правовые аспекты</h2>
                    <p>Использование программы для создания нагрузки на чужие серверы без разрешения их владельцев может квалифицироваться как нарушение законодательства Российской Федерации и других стран, включая статьи о компьютерных преступлениях.</p>
                    
                    <div class="responsibility-block">
                        <h3>Ответственность пользователя:</h3>
                        <ul>
                            <li>Получение разрешения на тестирование от владельца ресурса</li>
                            <li>Соблюдение законодательства страны пребывания</li>
                            <li>Использование программы исключительно в законных целях</li>
                            <li>Понимание последствий создания высокой нагрузки на серверы</li>
                        </ul>
                    </div>
                </div>

                <div class="license-section">
                    <h2>📄 Лицензия и авторство</h2>
                    <p><strong>Автор программы:</strong> Алексей Ульянов</p>
                    <p>Программа распространяется по лицензии MIT исключительно для ознакомительных целей и внутреннего тестирования нагрузки сайтов и серверов.</p>
                    <p>Автор не несет ответственности за использование программы в целях, нарушающих законодательство Российской Федерации и других стран.</p>
                </div>

                <div class="agreement-footer">
                    <p><strong>Используя данную программу, вы подтверждаете, что:</strong></p>
                    <ul>
                        <li>Ознакомились с данным соглашением и принимаете его условия</li>
                        <li>Понимаете принцип работы и возможные последствия использования</li>
                        <li>Обязуетесь использовать программу только в законных целях</li>
                        <li>Принимаете на себя полную ответственность за использование программы</li>
                    </ul>
                </div>
            </div>
        </div>
        
        <div class="disclaimer" style="text-align: center; padding: 20px; color: var(--text-muted); font-size: 12px; line-height: 1.6;">
            <strong>Автор программы:</strong> Алексей Ульянов | 
            <strong>Лицензия:</strong> MIT | 
            <strong>Версия:</strong> {{.Version}}
        </div>
    </div>

    <script>
        let updateInterval;
        
        // Инициализация темы при загрузке
        document.addEventListener('DOMContentLoaded', function() {
            initTheme();
            showCookieNotice();
        });

        function initTheme() {
            const savedTheme = localStorage.getItem('loadtester-theme') || 'light';
            const themeIcon = document.getElementById('themeIcon');
            const themeText = document.getElementById('themeText');
            
            if (savedTheme === 'dark') {
                document.documentElement.setAttribute('data-theme', 'dark');
                themeIcon.textContent = '☀️';
                themeText.textContent = 'Светлая';
            } else {
                document.documentElement.removeAttribute('data-theme');
                themeIcon.textContent = '🌙';
                themeText.textContent = 'Темная';
            }
        }

        function toggleTheme() {
            const currentTheme = document.documentElement.getAttribute('data-theme');
            const themeIcon = document.getElementById('themeIcon');
            const themeText = document.getElementById('themeText');
            
            if (currentTheme === 'dark') {
                document.documentElement.removeAttribute('data-theme');
                localStorage.setItem('loadtester-theme', 'light');
                themeIcon.textContent = '🌙';
                themeText.textContent = 'Темная';
            } else {
                document.documentElement.setAttribute('data-theme', 'dark');
                localStorage.setItem('loadtester-theme', 'dark');
                themeIcon.textContent = '☀️';
                themeText.textContent = 'Светлая';
            }
        }

        function showCookieNotice() {
            const cookieAccepted = localStorage.getItem('loadtester-cookies-accepted');
            if (!cookieAccepted) {
                setTimeout(() => {
                    document.getElementById('cookieNotice').classList.add('show');
                }, 1000);
            }
        }

        function acceptCookies() {
            localStorage.setItem('loadtester-cookies-accepted', 'true');
            document.getElementById('cookieNotice').classList.remove('show');
        }

        function showPage(pageId) {
            // Скрываем все страницы
            document.getElementById('mainPage').style.display = 'none';
            document.getElementById('instructionsPage').style.display = 'none';
            document.getElementById('agreementPage').style.display = 'none';
            
            // Показываем нужную страницу
            if (pageId === 'main') {
                document.getElementById('mainPage').style.display = 'block';
            } else if (pageId === 'instructions') {
                document.getElementById('instructionsPage').style.display = 'block';
            } else if (pageId === 'agreement') {
                document.getElementById('agreementPage').style.display = 'block';
            }
            
            // Обновляем активную ссылку в навигации
            document.querySelectorAll('.nav-link').forEach(link => {
                link.classList.remove('active');
            });
            
            if (pageId === 'main') {
                document.querySelector('.nav-link[onclick="showPage(\'main\')"]').classList.add('active');
            } else if (pageId === 'instructions') {
                document.querySelector('.nav-link[onclick="showPage(\'instructions\')"]').classList.add('active');
            } else if (pageId === 'agreement') {
                document.querySelector('.nav-link[onclick="showPage(\'agreement\')"]').classList.add('active');
            }
        }

        function toggleStartButton() {
            const checkbox = document.getElementById('agreementCheckbox');
            const startBtn = document.getElementById('startBtn');
            
            if (checkbox.checked) {
                startBtn.disabled = false;
            } else {
                startBtn.disabled = true;
            }
        }
        
        function startTest() {
            // Проверяем согласие с пользовательским соглашением
            const checkbox = document.getElementById('agreementCheckbox');
            if (!checkbox.checked) {
                alert('Необходимо принять пользовательское соглашение для запуска теста');
                return;
            }

            const form = document.getElementById('testForm');
            const formData = new FormData(form);
            const config = Object.fromEntries(formData);
            
            // Конвертируем числовые значения
            config.totalRequests = parseInt(config.totalRequests);
            config.maxConcurrency = parseInt(config.maxConcurrency);
            config.timeout = parseInt(config.timeout);
            config.maxMemory = parseFloat(config.maxMemory);
            
            fetch('/start', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify(config)
            })
            .then(response => response.json())
            .then(data => {
                if (data.success) {
                    document.getElementById('startBtn').disabled = true;
                    document.getElementById('stopBtn').disabled = false;
                    document.getElementById('stats').style.display = 'grid';
                    updateInterval = setInterval(updateStats, 500);
                } else {
                    alert('Ошибка: ' + data.error);
                }
            })
            .catch(error => {
                alert('Ошибка запроса: ' + error);
            });
        }
        
        function stopTest() {
            fetch('/stop', { method: 'POST' })
            .then(response => response.json())
            .then(data => {
                document.getElementById('stopBtn').disabled = true;
            });
        }
        
        function updateStats() {
            fetch('/stats')
            .then(response => response.json())
            .then(stats => {
                // Обновляем прогресс бар
                document.getElementById('progressFill').style.width = (stats.progress * 100) + '%';
                
                // Обновляем статус
                const statusEl = document.getElementById('status');
                if (stats.isRunning) {
                    statusEl.textContent = 'Выполняется: ' + stats.total + ' запросов | RPS: ' + stats.rps.toFixed(1);
                    statusEl.className = 'status running';
                } else {
                    const successRate = stats.total > 0 ? (stats.success / stats.total * 100) : 0;
                    let emoji = '🟢';
                    let text = 'Отличный результат!';
                    if (successRate <= 80) {
                        emoji = '🔴';
                        text = 'Есть проблемы';
                    } else if (successRate <= 95) {
                        emoji = '🟡';
                        text = 'Хороший результат';
                    }
                    statusEl.textContent = emoji + ' Тест завершен! ' + text + ' (' + successRate.toFixed(1) + '% успешных)';
                    statusEl.className = 'status completed';
                    
                    // Останавливаем обновления
                    clearInterval(updateInterval);
                    
                    // Включаем кнопку "Запустить тест" только если согласие принято
                    const checkbox = document.getElementById('agreementCheckbox');
                    document.getElementById('startBtn').disabled = !checkbox.checked;
                    document.getElementById('stopBtn').disabled = true;
                }
                
                // Обновляем статистику
                document.getElementById('totalStat').textContent = stats.total;
                document.getElementById('rpsStat').textContent = stats.rps.toFixed(1);
                document.getElementById('successStat').textContent = stats.success;
                document.getElementById('successPercent').textContent = stats.total > 0 ? (stats.success / stats.total * 100).toFixed(1) : '0';
                document.getElementById('failsStat').textContent = stats.fails;
                document.getElementById('status4xxStat').textContent = stats.status4xx;
                document.getElementById('status5xxStat').textContent = stats.status5xx;
                document.getElementById('errorsStat').textContent = stats.errors;
                document.getElementById('timeoutsStat').textContent = stats.timeouts;
                document.getElementById('connectionsStat').textContent = stats.connections;
            })
            .catch(error => {
                console.error('Ошибка получения статистики:', error);
            });
        }
    </script>
</body>
</html>
`
