const editButton = document.getElementById('editButton');
const overlay = document.getElementById('overlay');
const popupForm = document.getElementById('popupForm');

editButton.addEventListener('click', function() {
  overlay.style.display = 'block';
  popupForm.style.display = 'block';
  editButton.classList.add('active');
});

const closeButton = document.getElementById('closeButton');

closeButton.addEventListener('click', function() {
  overlay.style.display = 'none';
  popupForm.style.display = 'none';
  editButton.classList.remove('active');
});

const deleteButton = document.querySelector('.delete-button');
const overlay_delete = document.getElementById('overlay_delete');
const deleteFormContainer = document.querySelector('.delete-form-container');

deleteButton.addEventListener('click', function() {
    overlay_delete.style.display = 'block';
    deleteFormContainer.style.display = 'block';
    deleteButton.classList.add('active');
});

const deleteCloseButton = document.querySelector('.delete-close-button');

deleteCloseButton.addEventListener('click', function()  {
    overlay_delete.style.display = 'none';
    deleteFormContainer.style.display = 'none';
    deleteButton.classList.remove('active');
});

const buttons = document.querySelectorAll('.btn:not(.delete-button):not(#editButton)');

buttons.forEach(button => {
  button.addEventListener('click', function() {
    buttons.forEach(btn => btn.classList.remove('active'));
    this.classList.add('active');
  });
});
