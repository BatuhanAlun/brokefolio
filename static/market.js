
const tradeSymbolSelect = document.getElementById('trade-symbol');
const currentPriceSpan = document.getElementById('current-price');
const quantityInput = document.getElementById('quantity');
const totalPriceInput = document.getElementById('total-price');
const tradeMessageP = document.getElementById('trade-message');
const tradingViewChartContainer = document.getElementById('tradingViewChartContainer'); 
let currentPrice = 0;
function generateTradingViewWidgetHtml(symbol) {
    const newSymbol = getTradingViewSymbol(symbol);
    return `
        <div class="tradingview-widget-container__widget" style="height:calc(100% - 32px);width:100%"></div>
        <div class="tradingview-widget-copyright"><a href="https://www.tradingview.com/" rel="noopener nofollow" target="_blank"><span class="blue-text">Track all markets on TradingView</span></a></div>
        <script type="text/javascript" src="https://s3.tradingview.com/external-embedding/embed-widget-advanced-chart.js" async>
        {
            "allow_symbol_change": true,
            "calendar": false,
            "details": false,
            "hide_side_toolbar": true,
            "hide_top_toolbar": false,
            "hide_legend": false,
            "hide_volume": false,
            "hotlist": false,
            "interval": "D",
            "locale": "tr", // Set to Turkish
            "save_image": true,
            "style": "1",
            "symbol": "${newSymbol}", // THIS IS DYNAMICALLY INSERTED
            "theme": "dark",
            "timezone": "Etc/UTC",
            "backgroundColor": "rgba(26, 13, 47, 1)", // Matching your .chart-wrapper background-color #1A0D2F
            "gridColor": "rgba(46, 46, 46, 0.06)",
            "watchlist": [],
            "withdateranges": false,
            "compareSymbols": [],
            "studies": [],
            "autosize": true
        }
        </script>`;
}
function isCryptoSymbol(symbol) {
    const upperSymbol = symbol.toUpperCase();
    return upperSymbol.endsWith('USDT') || upperSymbol.endsWith('USD');
}

function getTradingViewSymbol(rawSymbol) {
        const symbolForBackend = rawSymbol.includes(':')
            ? rawSymbol.split(':')[1]
            : rawSymbol;
        if (symbolForBackend.endsWith('USD') || symbolForBackend.endsWith('USDT') ) {
            
            const fullTradingViewSymbol = `BINANCE:${symbolForBackend}`;
            return fullTradingViewSymbol
        } else {
            const fullTradingViewSymbol = symbolForBackend;
            return fullTradingViewSymbol
        }
}

async function fetchCurrentPrice(symbolFromDropdown) {
    tradeMessageP.textContent = '';
    currentPriceSpan.textContent = 'Yükleniyor...';
    try {

        const symbolForBackend = symbolFromDropdown.includes(':')
            ? symbolFromDropdown.split(':')[1]
            : symbolFromDropdown;
        let apiEndpoint;

        if (symbolForBackend.endsWith('USD')) {
            apiEndpoint = `/api/crypto-price?symbol=${symbolForBackend}`;
        } else {
            apiEndpoint = `/api/crypto-price?symbol=${symbolForBackend}`;
        }
        const response = await fetch(apiEndpoint);
        if (!response.ok) {
            const errorText = await response.text();
            throw new Error(`HTTP error! Status: ${response.status} - ${errorText}`);
        }
        const data = await response.json();
        if (data.price) {
            currentPrice = parseFloat(data.price);
            currentPriceSpan.textContent = `${currentPrice.toFixed(4)} USD`;
            calculateTotalPrice();
        } else {
            currentPriceSpan.textContent = 'Fiyat bulunamadı.';
            currentPrice = 0;
            totalPriceInput.value = '';
        }
    } catch (error) {
        console.error('Fiyat yüklenirken hata oluştu:', error);
        currentPriceSpan.textContent = 'Hata!';
        currentPrice = 0;
        totalPriceInput.value = '';
    }
}
function calculateTotalPrice() {
    const quantity = parseFloat(quantityInput.value);
    if (!isNaN(quantity) && currentPrice > 0) {
        const total = quantity * currentPrice;
        totalPriceInput.value = total.toFixed(4);
    } else {
        totalPriceInput.value = '';
    }
}
async function submitTrade(tradeType) {
    tradeMessageP.textContent = 'İşlem yapılıyor...';
    
    const symbolFromDropdown = tradeSymbolSelect.value; 
    const symbolForBackend = symbolFromDropdown.includes(':')
        ? symbolFromDropdown.split(':')[1]
        : symbolFromDropdown;
    const quantity = parseFloat(quantityInput.value);
    const price = currentPrice;
    if (isNaN(quantity) || quantity <= 0) {
        tradeMessageP.style.color = '#FF4C4C';
        tradeMessageP.textContent = 'Lütfen geçerli bir miktar girin.';
        return;
    }
    if (price <= 0) {
        tradeMessageP.style.color = '#FF4C4C';
        tradeMessageP.textContent = 'Güncel fiyat mevcut değil, lütfen bekleyin.';
        return;
    }
    const tradeData = {
        symbol: symbolForBackend,
        quantity: quantity,
        price: price,
        type: tradeType
    };
    try {
        const response = await fetch('/api/trade', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify(tradeData)
        });
        const result = await response.json();
        if (response.ok) {
            tradeMessageP.style.color = '#4CAF50';
            tradeMessageP.textContent = result.message;
            quantityInput.value = '';
            totalPriceInput.value = '';
        } else {
            tradeMessageP.style.color = '#FF4C4C';
            tradeMessageP.textContent = result.error || 'İşlem başarısız oldu.';
        }
    }
    catch (error) {
        console.error('İşlem gönderilirken hata oluştu:', error);
        tradeMessageP.style.color = '#FF4C4C';
        tradeMessageP.textContent = 'Bağlantı hatası: İşlem gönderilemedi.';
    }
}

async function fetchCombinedNews() {
    const newsFeedDiv = document.getElementById('crypto-news-feed');
    newsFeedDiv.innerHTML = '<p style="color: #ccc; text-align: center; width: 100%;">Haberler yükleniyor...</p>';
    try {
        const response = await fetch('/api/news');
        if (!response.ok) {
            throw new Error(`HTTP error! Status: ${response.status}`);
        }
        const data = await response.json();
        newsFeedDiv.innerHTML = '';
        if (data && data.results && data.results.length > 0) {
            data.results.forEach(news => {
                const newsItem = document.createElement('div');
                newsItem.classList.add('news-item');
                const dateText = news.date;
                newsItem.innerHTML = `
                    <h4>${news.title}</h4>
                    <p>${news.source} - <span class="source-date">${dateText}</span></p>
                    <a href="${news.url}" target="_blank" rel="noopener">Devamını Oku <i class="fas fa-external-link-alt"></i></a>
                `;
                newsFeedDiv.appendChild(newsItem);
            });
        } else {
            newsFeedDiv.innerHTML = '<p style="color: #ccc; text-align: center; width: 100%;">Şu an için haber bulunmuyor.</p>';
        }
    } catch (error) {
        console.error("Haberler yüklenirken hata oluştu:", error);
        newsFeedDiv.innerHTML = '<p style="color: #dc3545; text-align: center; width: 100%;">Haberler yüklenemedi. Lütfen daha sonra tekrar deneyin.</p>';
    }
}

tradeSymbolSelect.addEventListener('change', () => {
    const selectedSymbol = tradeSymbolSelect.value;
    
    fetchCurrentPrice(selectedSymbol);

    tradingViewChartContainer.innerHTML = generateTradingViewWidgetHtml(selectedSymbol);
    console.log(`TradingView widget reloaded with symbol: ${selectedSymbol}`);
});
quantityInput.addEventListener('input', calculateTotalPrice);

document.addEventListener('DOMContentLoaded', () => {
    fetchCurrentPrice(tradeSymbolSelect.value); 
    tradingViewChartContainer.innerHTML = generateTradingViewWidgetHtml(tradeSymbolSelect.value);
    fetchCombinedNews();
});
