document.addEventListener('DOMContentLoaded', () => {
    // DOM Elements
    const searchInput = document.getElementById('search-input');
    const searchButton = document.getElementById('search-button');
    const advancedSearchToggle = document.getElementById('advanced-search-toggle');
    const advancedFilters = document.getElementById('advanced-filters');
    const applyFiltersButton = document.getElementById('apply-filters');
    const clearFiltersButton = document.getElementById('clear-filters');
    const resultsContainer = document.getElementById('results');
    const resultCount = document.getElementById('result-count');
    const prevPageButton = document.getElementById('prev-page');
    const nextPageButton = document.getElementById('next-page');
    const pageInfo = document.getElementById('page-info');
    const modal = document.getElementById('log-details-modal');
    const closeModal = document.querySelector('.close');
    const logJson = document.getElementById('log-json');

    // State
    let currentPage = 1;
    let totalPages = 1;
    let currentLogs = [];
    const pageSize = 10;

    // Toggle advanced filters
    advancedSearchToggle.addEventListener('click', () => {
        advancedFilters.classList.toggle('show');
        const icon = advancedSearchToggle.querySelector('i');
        if (advancedFilters.classList.contains('show')) {
            icon.className = 'fas fa-chevron-up';
        } else {
            icon.className = 'fas fa-chevron-down';
        }
    });

    // Search on Enter key
    searchInput.addEventListener('keyup', (e) => {
        if (e.key === 'Enter') {
            searchLogs();
        }
    });

    // Search button click
    searchButton.addEventListener('click', searchLogs);

    // Apply filters button click
    applyFiltersButton.addEventListener('click', searchLogs);

    // Clear filters button click
    clearFiltersButton.addEventListener('click', clearFilters);

    // Pagination
    prevPageButton.addEventListener('click', () => {
        if (currentPage > 1) {
            currentPage--;
            displayLogs();
        }
    });

    nextPageButton.addEventListener('click', () => {
        if (currentPage < totalPages) {
            currentPage++;
            displayLogs();
        }
    });

    // Close modal
    closeModal.addEventListener('click', () => {
        modal.style.display = 'none';
    });

    // Close modal when clicking outside
    window.addEventListener('click', (e) => {
        if (e.target === modal) {
            modal.style.display = 'none';
        }
    });

    // Format timestamp
    function formatTimestamp(timestamp) {
        const date = new Date(timestamp);
        return date.toLocaleString();
    }

    // Clear all filters
    function clearFilters() {
        document.getElementById('level').value = '';
        document.getElementById('resourceId').value = '';
        document.getElementById('traceId').value = '';
        document.getElementById('spanId').value = '';
        document.getElementById('commit').value = '';
        document.getElementById('parentResourceId').value = '';
        document.getElementById('startTime').value = '';
        document.getElementById('endTime').value = '';
        document.getElementById('regex').value = '';
        document.getElementById('message').value = '';
        searchInput.value = '';
    }

    // Build query parameters
    function buildQueryParams() {
        const params = new URLSearchParams();
        
        // Full-text search
        const searchValue = searchInput.value.trim();
        if (searchValue) {
            params.append('search', searchValue);
        }
        
        // Level
        const level = document.getElementById('level').value;
        if (level) {
            params.append('level', level);
        }
        
        // ResourceId
        const resourceId = document.getElementById('resourceId').value.trim();
        if (resourceId) {
            params.append('resourceId', resourceId);
        }
        
        // TraceId
        const traceId = document.getElementById('traceId').value.trim();
        if (traceId) {
            params.append('traceId', traceId);
        }
        
        // SpanId
        const spanId = document.getElementById('spanId').value.trim();
        if (spanId) {
            params.append('spanId', spanId);
        }
        
        // Commit
        const commit = document.getElementById('commit').value.trim();
        if (commit) {
            params.append('commit', commit);
        }
        
        // ParentResourceId
        const parentResourceId = document.getElementById('parentResourceId').value.trim();
        if (parentResourceId) {
            params.append('parentResourceId', parentResourceId);
        }
        
        // Start Time
        const startTime = document.getElementById('startTime').value;
        if (startTime) {
            params.append('startTime', new Date(startTime).toISOString());
        }
        
        // End Time
        const endTime = document.getElementById('endTime').value;
        if (endTime) {
            params.append('endTime', new Date(endTime).toISOString());
        }
        
        // Regex Pattern
        const regex = document.getElementById('regex').value.trim();
        if (regex) {
            params.append('regex', regex);
        }
        
        // Message
        const message = document.getElementById('message').value.trim();
        if (message) {
            params.append('message', message);
        }
        
        // Pagination
        params.append('page', currentPage);
        params.append('limit', pageSize);
        
        return params;
    }

    // Search logs
    async function searchLogs() {
        try {
            currentPage = 1;
            const params = buildQueryParams();
            const response = await fetch(`/logs?${params.toString()}`);
            
            if (!response.ok) {
                throw new Error(`HTTP error! Status: ${response.status}`);
            }
            
            const data = await response.json();
            currentLogs = data.logs || [];
            
            // Calculate total pages
            const totalLogs = data.count || 0;
            totalPages = Math.ceil(totalLogs / pageSize);
            
            displayLogs();
        } catch (error) {
            console.error('Error searching logs:', error);
            resultsContainer.innerHTML = `
                <div class="no-results">
                    <i class="fas fa-exclamation-circle fa-3x"></i>
                    <p>Error fetching logs: ${error.message}</p>
                </div>
            `;
            resultCount.textContent = '(0)';
            updatePaginationControls();
        }
    }

    // Display logs
    function displayLogs() {
        if (currentLogs.length === 0) {
            resultsContainer.innerHTML = `
                <div class="no-results">
                    <i class="fas fa-search fa-3x"></i>
                    <p>No logs found. Try adjusting your search criteria.</p>
                </div>
            `;
            resultCount.textContent = '(0)';
        } else {
            let html = '';
            
            currentLogs.forEach(log => {
                html += `
                    <div class="log-item" data-log='${JSON.stringify(log)}'>
                        <div class="log-header">
                            <span class="log-level ${log.level.toLowerCase()}">${log.level}</span>
                            <span class="log-timestamp">${formatTimestamp(log.timestamp)}</span>
                        </div>
                        <div class="log-message">${log.message}</div>
                        <div class="log-details">
                            <div class="log-detail">
                                <i class="fas fa-server"></i>
                                <span>${log.resourceId}</span>
                            </div>
                            <div class="log-detail">
                                <i class="fas fa-fingerprint"></i>
                                <span>${log.traceId}</span>
                            </div>
                            <div class="log-detail">
                                <i class="fas fa-code-branch"></i>
                                <span>${log.commit}</span>
                            </div>
                        </div>
                    </div>
                `;
            });
            
            resultsContainer.innerHTML = html;
            resultCount.textContent = `(${currentLogs.length})`;
            
            // Add click event to log items
            document.querySelectorAll('.log-item').forEach(item => {
                item.addEventListener('click', () => {
                    const logData = JSON.parse(item.getAttribute('data-log'));
                    showLogDetails(logData);
                });
            });
        }
        
        // Update pagination
        pageInfo.textContent = `Page ${currentPage}`;
        updatePaginationControls();
    }

    // Update pagination controls
    function updatePaginationControls() {
        prevPageButton.disabled = currentPage <= 1;
        nextPageButton.disabled = currentPage >= totalPages;
    }

    // Show log details in modal
    function showLogDetails(log) {
        logJson.textContent = JSON.stringify(log, null, 2);
        modal.style.display = 'block';
    }

    // Initial search on page load
    searchLogs();
}); 