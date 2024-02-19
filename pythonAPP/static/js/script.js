document.addEventListener('DOMContentLoaded', function() {
  const createUserBtn = document.getElementById('createUserBtn');
  const newProductBtn = document.getElementById('newProduct');
  const editProductBtns = document.querySelectorAll('.editProduct');
  const addClientStockBtn = document.getElementById('addClientStockBtn');
  const reduceButton = document.getElementById('reduceClientStock');

  const overlay_create = document.getElementById('overlay_create');
  const overlay_product = document.getElementById('overlay_product');
  const overlay_productedits = document.querySelectorAll('.overlay');
  const overlay_addClientStock = document.getElementById('addClientStockOverlay');
  const overlay_recudeClientStock = document.getElementById('reduceClientStockOverlay');

  const popupFormCreate = document.getElementById('popupFormCreate');
  const popupFormProduct = document.getElementById('popupFormProduct');
  const popupFormProductedit = document.querySelectorAll('.form-container')
  const popupFormAddClientStock = document.getElementById('popupformAddClientStock');
  const popupFormReduceClientStock = document.getElementById('popupformreduceClientStock');

  const deleteButton = document.querySelector('.delete-button');
  const overlay_delete = document.getElementById('overlay_delete');
  const deleteFormContainer = document.querySelector('.delete-form-container');
  const productCloseButtonsedit = document.querySelectorAll('#productCloseButtonedit');
  // const closeAddClientStockBtn = document.getElementsByName('newCloseBtnnn')
  const reduceClientStockCloseBtn = document.getElementById('close-buttonReduce');

  const editButton = document.getElementById('editButton');
  const overlay_edit = document.getElementById('overlayedit');
  const editFormContainer = document.getElementById('popupFormEdit');

  if (reduceButton) {
    reduceButton.addEventListener('click', function() {
      overlay_recudeClientStock.style.display = 'block';
      popupFormReduceClientStock.style.display = 'block';
    })
  }

  if (reduceClientStockCloseBtn) {
    overlay_recudeClientStock.style.display = 'none';
    popupFormReduceClientStock.style.display = 'none';
  }

  if (addClientStockBtn) {
    addClientStockBtn.addEventListener('click', function() {
      overlay_addClientStock.style.display = 'block';
      popupFormAddClientStock.style.display = 'block';
    })
  }

//   if (closeAddClientStockBtn) {
//     closeAddClientStockBtn.addEventListener('click', function() {
//     overlay_addClientStock.style.display = 'none';
//     popupFormAddClientStock.style.display = 'none';
//     console.log('Close button clicked');
// })
// }


  editProductBtns.forEach((editProductBtn, index) => {
    editProductBtn.addEventListener('click', function() {
        overlay_productedits[index+1].style.display = 'block';
        popupFormProductedit[index+1].style.display = 'block';
    });
});

productCloseButtonsedit.forEach((productCloseButtonedit, index) => {
  productCloseButtonedit.addEventListener('click', function() {
      overlay_productedits[index+1].style.display = 'none';
      // popupFormProductedits[index+1].style.display = 'none';
  });
});

// productCloseButtonsedit.forEach((productCloseButtonedit, index) => {
//   productCloseButtonedit.addEventListener('click', function() {
      // overlay_productedits[index+1].style.display = 'none';
      // popupFormProductedit[index+1].style.display = 'none';
//       // ... your other logic ...
//   });
// });

  if (editButton) {
    editButton.addEventListener('click', function() {
      overlay_edit.style.display = 'block';
      editFormContainer.style.display = 'block';
      // createUserBtn.classList.add('active');
    })
  }

  const editCloseButton = document.getElementById('closeeditButton');

  if (editCloseButton) {
      editCloseButton.addEventListener('click', function() {
      overlay_edit.style.display = 'none';
      editFormContainer.style.display = 'none';
      // createUserBtn.classList.add('active');
    })
  }

  const addButtons = document.querySelectorAll('.btn-add-stock');
  // const add_overlay = document.getElementById("overlay_add");
  // const addFormContainer = document.getElementById("add-form-container");
  const closeButton = document.getElementById('addingBtn');


  function openModal(data) {
    console.log("Opening modal");
    const overlay = document.getElementById('overlay_add');
    const formContainer = document.getElementById('add-form-container');
    console.log("Overlay and formContainer:", overlay, formContainer);

    const productNameElement = document.getElementById('productName');
    const productIDElement = document.getElementById('productID');
    console.log("productNameElement and productIDElement:", productNameElement, productIDElement);

    if (productNameElement && productIDElement) {
        productNameElement.innerText = data.productName;
        productIDElement.value = data.productID;
        overlay.style.display = 'block';
        formContainer.style.display = 'block';
        console.log("Modal opened successfully");
    } else {
        console.log("Unable to open modal - elements not found");
    }
}

function closeModal() {
  const overlay = document.getElementById('overlay_add');
  overlay.style.display = 'none';
  formContainer.style.display= 'none';
}

addButtons.forEach(button => {
  button.addEventListener('click', function (event) {
      console.log("Add Stock button clicked");
      event.preventDefault();
      const data = {
          productName: button.getAttribute('data-product-name'),
          productID: button.getAttribute('data-product-id')
      };
      openModal(data);
  });
});

if (closeButton) {
  closeButton.addEventListener('click', function () {
    console.log("Close button clicked");
    closeModal();
  });
}

// if (addCloseBtn) {
//   addCloseBtn.addEventListener('click', closeModal);
// }
  // if (addBtn) {
  //   addBtn.addEventListener('click', function() {
  //     add_overlay.style.display = 'block';
  //     addFormContainer.style.display = 'block';
  //   })
  // }


  // if (addCloseBtn) {
  //   addCloseBtn.addEventListener('click', function() {
  //     add_overlay.style.display = 'none';
  //     addFormContainer.style.display = 'none';
  //   })
  // }

  if (createUserBtn) {
      createUserBtn.addEventListener('click', function() {
          overlay_create.style.display = 'block';
          popupFormCreate.style.display = 'block';
          // createUserBtn.classList.add('active');
      });
  }

  if (deleteButton) {
      deleteButton.addEventListener('click', function() {
        overlay_delete.style.display = 'block';
        deleteFormContainer.style.display = 'block';
        // deleteButton.classList.add('active');
      });
  }

  const deleteCloseButton = document.getElementById('closingBtn');

  if (deleteCloseButton) {
    deleteCloseButton.addEventListener('click', function()  {
      overlay_delete.style.display = 'none';
      deleteFormContainer.style.display = 'none';
      // deleteButton.classList.remove('active');
    });
  }

  if (newProductBtn) {
      newProductBtn.addEventListener('click', function() {
          overlay_product.style.display = 'block';
          popupFormProduct.style.display = 'block';
          // newProductBtn.classList.add('active');
      });
  }

  const createCloseButton = document.getElementById('createCloseButton');
  const productCloseButton = document.getElementById('productCloseButton');

  if (createCloseButton) {
      createCloseButton.addEventListener('click', function() {
          overlay_create.style.display = 'none';
          popupFormCreate.style.display = 'none';
          // createUserBtn.classList.remove('active');
      });
  }

  if (productCloseButton) {
      productCloseButton.addEventListener('click', function() {
          overlay_product.style.display = 'none';
          popupFormProduct.style.display = 'none';
          // newProductBtn.classList.remove('active');
      });
  }

  const buttons = document.querySelectorAll('.btn:not(.delete-button):not(#editButton)');

  buttons.forEach(button => {
    button.addEventListener('click', function() {
    buttons.forEach(btn => btn.classList.remove('active'));
    this.classList.add('active');
  });
});
});


// const createUserBtn = document.getElementById('createUserBtn');
// const overlay_create = document.getElementById('overlay_create');
// const popupFormCreate = document.getElementById('popupFormCreate');

// createUserBtn.addEventListener('click', function() {
//   overlay_create.style.display = 'block';
//   popupFormCreate.style.display = 'block';
//   // createUserBtn.classList.add('active');
// });

// const createCloseButton = document.getElementById('createCloseButton');

// createCloseButton.addEventListener('click', function() {
//   overlay_create.style.display = 'none';
//   popupFormCreate.style.display = 'none';
//   // createUserBtn.classList.remove('active');
// });

// const createProductBtn = document.getElementsByName('newProduct');
// const overlay_product = document.getElementById('overlay_product');
// const popupFormProduct = document.getElementById('popupFormProduct');

// createProductBtn.addEventListener('click', function() {
//   overlay_product.style.display = 'block';
//   popupFormProduct.style.display = 'block';
//   // createProductBtn.classList.add('active');
// });

// const productCloseButton = document.getElementById('productCloseButton');

// productCloseButton.addEventListener('click', function() {
//   overlay_product.style.display = 'none';
//   popupFormProduct.style.display = 'none';
//   // createProductBtn.classList.remove('active');
// });
