function toggleFilters() {
    const filterContainer = document.getElementById('filter-container');
    if (filterContainer.style.display === 'none' || filterContainer.style.display === '') {
        filterContainer.style.display = 'grid';
    } else {
        filterContainer.style.display = 'none';
    }
}

function toggleMoreOptions() {
    const filterContainer = document.getElementById('moreOptions');
    if (filterContainer.style.display === 'none' || filterContainer.style.display === '') {
        filterContainer.style.display = 'grid';
    } else {
        filterContainer.style.display = 'none';
    }
}