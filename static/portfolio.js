const totalPortfolioValueSpan = document.getElementById('total-portfolio-value');
        const totalPnlSpan = document.getElementById('total-pnl');
        const portfolioHoldingsBody = document.getElementById('portfolio-holdings-body');
        const portfolioMessageP = document.getElementById('portfolio-message');
        const transactionsBody = document.getElementById('transactions-body');
        const transactionsMessageP = document.getElementById('transactions-message');

        const sellAssetSelect = document.getElementById('sell-asset-select');
        const selectedAssetQuantitySpan = document.getElementById('selected-asset-quantity');
        const selectedAssetPriceSpan = document.getElementById('selected-asset-price');
        const sellQuantityInput = document.getElementById('sell-quantity-input');
        const sellTotalValueP = document.getElementById('sell-total-value');
        const executeSellButton = document.getElementById('execute-sell-button');
        const sellSectionMessageP = document.getElementById('sell-section-message');


        let currentHoldings = [];

        function formatCurrency(value) {
            return value.toLocaleString('en-US', { style: 'currency', currency: 'USD' });
        }

        function formatPercentage(value) {
            return value.toFixed(2) + '%';
        }

        function formatTimestamp(timestamp) {
            const date = new Date(timestamp);
            return date.toLocaleString();
        }

        function populateSellAssetDropdown(holdings) {
            sellAssetSelect.innerHTML = '<option value="">-- Varlık Seç --</option>';
            if (holdings && holdings.length > 0) {
                holdings.forEach(holding => {
                    const option = document.createElement('option');
                    option.value = holding.symbol;
                    option.textContent = `${holding.symbol} (Mevcut: ${holding.quantity.toFixed(6)})`;
                    option.dataset.quantity = holding.quantity;
                    option.dataset.price = holding.currentPrice;
                    sellAssetSelect.appendChild(option);
                });
            }
        }

        sellAssetSelect.addEventListener('change', () => {
            const selectedOption = sellAssetSelect.options[sellAssetSelect.selectedIndex];
            if (selectedOption && selectedOption.value) {
                const quantity = parseFloat(selectedOption.dataset.quantity);
                const price = parseFloat(selectedOption.dataset.price);

                selectedAssetQuantitySpan.textContent = quantity.toFixed(6);
                selectedAssetPriceSpan.textContent = formatCurrency(price);
                sellQuantityInput.value = '';
                sellTotalValueP.textContent = 'Toplam Değer: $0.00';
                sellSectionMessageP.textContent = '';

                sellAssetSelect.dataset.selectedQuantity = quantity;
                sellAssetSelect.dataset.selectedPrice = price;
            } else {

                selectedAssetQuantitySpan.textContent = '0';
                selectedAssetPriceSpan.textContent = formatCurrency(0);
                sellQuantityInput.value = '';
                sellTotalValueP.textContent = 'Toplam Değer: $0.00';
                sellSectionMessageP.textContent = '';
            }
        });

        sellQuantityInput.addEventListener('input', () => {
            const quantity = parseFloat(sellQuantityInput.value);
            const selectedPrice = parseFloat(sellAssetSelect.dataset.selectedPrice);

            if (!isNaN(quantity) && quantity > 0 && !isNaN(selectedPrice)) {
                sellTotalValueP.textContent = `Toplam Değer: ${formatCurrency(quantity * selectedPrice)}`;
            } else {
                sellTotalValueP.textContent = 'Toplam Değer: $0.00';
            }
            sellSectionMessageP.textContent = '';
        });

        executeSellButton.addEventListener('click', async () => {
            sellSectionMessageP.textContent = 'Satış işlemi yapılıyor...';
            sellSectionMessageP.style.color = '#FFD700';

            const symbol = sellAssetSelect.value;
            const quantityToSell = parseFloat(sellQuantityInput.value);
            const availableQuantity = parseFloat(sellAssetSelect.dataset.selectedQuantity);
            const price = parseFloat(sellAssetSelect.dataset.selectedPrice);

            if (!symbol) {
                sellSectionMessageP.style.color = '#FF4C4C';
                sellSectionMessageP.textContent = 'Lütfen satılacak bir varlık seçin.';
                return;
            }
            if (isNaN(quantityToSell) || quantityToSell <= 0) {
                sellSectionMessageP.style.color = '#FF4C4C';
                sellSectionMessageP.textContent = 'Lütfen geçerli bir miktar girin.';
                return;
            }
            if (isNaN(price) || price <= 0) {
                sellSectionMessageP.style.color = '#FF4C4C';
                sellSectionMessageP.textContent = 'Seçilen varlığın güncel fiyatı mevcut değil.';
                return;
            }
            if (quantityToSell > availableQuantity) {
                sellSectionMessageP.style.color = '#FF4C4C';
                sellSectionMessageP.textContent = `Yeterli ${symbol} miktarınız bulunmamaktadır. Mevcut: ${availableQuantity.toFixed(6)}`;
                return;
            }

            const tradeData = {
                symbol: symbol,
                quantity: quantityToSell,
                price: price,
                type: "SELL"
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
                    sellSectionMessageP.style.color = '#4CAF50';
                    sellSectionMessageP.textContent = result.message || 'Varlık başarıyla satıldı!';
                    // Refresh all data after successful sell
                    fetchPortfolioData();
                    fetchTransactionsData();
                    // Reset sell form
                    sellAssetSelect.value = '';
                    selectedAssetQuantitySpan.textContent = '0';
                    selectedAssetPriceSpan.textContent = formatCurrency(0);
                    sellQuantityInput.value = '';
                    sellTotalValueP.textContent = 'Toplam Değer: $0.00';
                } else {
                    sellSectionMessageP.style.color = '#FF4C4C';
                    sellSectionMessageP.textContent = result.error || 'Satış işlemi başarısız oldu.';
                }
            } catch (error) {
                console.error('Satış işlemi gönderilirken hata oluştu:', error);
                sellSectionMessageP.style.color = '#FF4C4C';
                sellSectionMessageP.textContent = 'Bağlantı hatası: Satış işlemi gönderilemedi.';
            }
        });

        async function fetchPortfolioData() {
            portfolioHoldingsBody.innerHTML = '<tr><td colspan="7" style="text-align: center; color: #ccc;">Portföy verileri yükleniyor...</td></tr>';
            portfolioMessageP.textContent = '';
            totalPortfolioValueSpan.textContent = 'Yükleniyor...';
            totalPnlSpan.textContent = 'Yükleniyor...';
            totalPnlSpan.classList.remove('positive', 'negative');

            try {
                const response = await fetch('/api/portfolio');
                if (!response.ok) {
                    throw new Error(`HTTP error! Status: ${response.status}`);
                }
                const data = await response.json();
                console.log("Portfolio Data:", data);

                if (data.error) {
                    portfolioMessageP.style.color = '#FF4C4C';
                    portfolioMessageP.textContent = data.error;
                    portfolioHoldingsBody.innerHTML = '<tr><td colspan="7" style="text-align: center; color: #dc3545;">Portföy yüklenirken hata oluştu.</td></tr>';
                    totalPortfolioValueSpan.textContent = formatCurrency(0);
                    totalPnlSpan.textContent = formatCurrency(0);
                    populateSellAssetDropdown([]);
                    return;
                }

                if (!data.holdings || data.holdings.length === 0) {
                    portfolioHoldingsBody.innerHTML = '<tr><td colspan="7" style="text-align: center; color: #ccc;">Portföyünüzde henüz bir varlık bulunmamaktadır.</td></tr>';
                    totalPortfolioValueSpan.textContent = formatCurrency(0);
                    totalPnlSpan.textContent = formatCurrency(0);
                    totalPnlSpan.classList.remove('positive', 'negative');
                    populateSellAssetDropdown([]);
                    return;
                }


                currentHoldings = data.holdings;
                populateSellAssetDropdown(currentHoldings);

                let totalPortfolioValue = 0;
                let totalUnrealizedPnl = 0;

                portfolioHoldingsBody.innerHTML = '';

                currentHoldings.forEach(holding => {
                    const row = document.createElement('tr');
                    const currentValue = holding.quantity * holding.currentPrice;
                    const pnl = currentValue - (holding.quantity * holding.averageBuyPrice);
                    const pnlPercentage = (holding.averageBuyPrice > 0) ? (pnl / (holding.quantity * holding.averageBuyPrice)) * 100 : 0;

                    totalPortfolioValue += currentValue;
                    totalUnrealizedPnl += pnl;

                    const pnlClass = pnl >= 0 ? 'value-positive' : 'value-negative';

                    row.innerHTML = `
                        <td>${holding.symbol}</td>
                        <td>${holding.quantity.toFixed(6)}</td>
                        <td>${formatCurrency(holding.averageBuyPrice)}</td>
                        <td>${formatCurrency(holding.currentPrice)}</td>
                        <td>${formatCurrency(currentValue)}</td>
                        <td class="${pnlClass}">${formatCurrency(pnl)}</td>
                        <td class="${pnlClass}">${formatPercentage(pnlPercentage)}</td>
                    `;
                    portfolioHoldingsBody.appendChild(row);
                });

                totalPortfolioValueSpan.textContent = formatCurrency(totalPortfolioValue);
                totalPnlSpan.textContent = formatCurrency(totalUnrealizedPnl);
                totalPnlSpan.classList.add(totalUnrealizedPnl >= 0 ? 'positive' : 'negative');

            } catch (error) {
                console.error("Portföy verileri yüklenirken hata oluştu:", error);
                portfolioMessageP.style.color = '#dc3545';
                portfolioMessageP.textContent = 'Portföy verileri yüklenemedi. Lütfen daha sonra tekrar deneyin.';
                portfolioHoldingsBody.innerHTML = '<tr><td colspan="7" style="text-align: center; color: #dc3545;">Veri yüklenirken hata oluştu.</td></tr>';
                totalPortfolioValueSpan.textContent = formatCurrency(0);
                totalPnlSpan.textContent = formatCurrency(0);
                totalPnlSpan.classList.remove('positive', 'negative');
                populateSellAssetDropdown([]);
            }
        }

        async function fetchTransactionsData() {
            transactionsBody.innerHTML = '<tr><td colspan="5" style="text-align: center; color: #ccc;">İşlem geçmişi yükleniyor...</td></tr>';
            transactionsMessageP.textContent = '';

            try {
                const response = await fetch('/api/transactions');
                if (!response.ok) {
                    throw new Error(`HTTP error! Status: ${response.status}`);
                }
                const data = await response.json();
                console.log("Transactions Data:", data);

                if (data.error) {
                    transactionsMessageP.style.color = '#FF4C4C';
                    transactionsMessageP.textContent = data.error;
                    transactionsBody.innerHTML = '<tr><td colspan="5" style="text-align: center; color: #dc3545;">İşlem geçmişi yüklenirken hata oluştu.</td></tr>';
                    return;
                }

                if (!data.transactions || data.transactions.length === 0) {
                    transactionsBody.innerHTML = '<tr><td colspan="5" style="text-align: center; color: #ccc;">Henüz bir işlem bulunmamaktadır.</td></tr>';
                    return;
                }

                transactionsBody.innerHTML = '';
                data.transactions.forEach(transaction => {
                    const row = document.createElement('tr');

                    let typeDisplay = transaction.type.trim().toUpperCase();
                    let colorClass = '';
                    
                    if (typeDisplay === 'BUY') {
                        colorClass = 'buy';
                    } else if (typeDisplay === 'SELL') {
                        colorClass = 'sell';
                    } else {
                        colorClass = 'transaction-neutral';
                    }
                    row.innerHTML = `
                        <td>${formatTimestamp(transaction.timestamp)}</td>
                        <td>${transaction.symbol}</td>
                        <td class="transaction-type ${colorClass}">${typeDisplay}</td>
                        <td>${transaction.quantity.toFixed(6)}</td>
                        <td>${formatCurrency(transaction.price)}</td>
                    `;
                    transactionsBody.appendChild(row);
                });

            } catch (error) {
                console.error("İşlem geçmişi yüklenirken hata oluştu:", error);
                transactionsMessageP.style.color = '#dc3545';
                transactionsMessageP.textContent = 'İşlem geçmişi yüklenemedi. Lütfen daha sonra tekrar deneyin.';
                transactionsBody.innerHTML = '<tr><td colspan="5" style="text-align: center; color: #dc3545;">Veri yüklenirken hata oluştu.</td></tr>';
            }
        }


        document.addEventListener('DOMContentLoaded', () => {
            fetchPortfolioData();
            fetchTransactionsData();
        });