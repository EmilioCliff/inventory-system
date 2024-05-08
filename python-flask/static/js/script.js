document.addEventListener('DOMContentLoaded', function() {
  const createUserBtn = document.getElementById('createUserBtn');
  const newProductBtn = document.getElementById('newProduct');
  const editProductBtns = document.querySelectorAll('.editProduct');
  const addClientStockBtn = document.getElementById('addClientStockBtn');
  const reduceButton = document.querySelector('#reduceClientStock');
  const waitViewButton = document.querySelector('#initializeSTK');
  const viewReceiptBtns = document.querySelectorAll('#viewReceipt');
  const viewInvoiceBtns = document.querySelectorAll('#viewInvoice');
  const viewTransactionBtns = document.querySelectorAll('#viewTransaction');
  const viewUserTransactionBtns = document.querySelectorAll('#viewUserTransaction');
  const viewDebtBtns = document.querySelectorAll('#viewDebt');

  const overlay_create = document.getElementById('overlay_create');
  const overlay_product = document.getElementById('overlay_product');
  const overlay_productedits = document.querySelectorAll('.overlay');
  const overlay_addClientStock = document.getElementById('addClientStockOverlay');
  const overlay_recudeClientStock = document.getElementById('reduceClientStockOverlay');
  const overlay_Wait = document.getElementById('waitOverlay');
  const overlay_requestClientStock = document.getElementById('requestClientStockOverlay');
  const overlay_viewReceipts = document.querySelectorAll('#overlay_receiptview');
  const overlay_viewInvoices = document.querySelectorAll('#overlay_invoiceview');
  const overlay_viewTransaction = document.querySelectorAll('#overlay_transactionview');
  const overlay_viewUserTransaction = document.querySelectorAll('#overlay_transactionuserview');
  const overlay_viewDebt = document.querySelectorAll('#overlay_debtview');

  const popupFormCreate = document.getElementById('popupFormCreate');
  const popupFormProduct = document.getElementById('popupFormProduct');
  const popupFormProductedit = document.querySelectorAll('.form-container')
  const popupFormAddClientStock = document.getElementById('popupformAddClientStock');
  const popupFormReduceClientStock = document.getElementById('popupformreduceClientStock');
  const popupFormWait = document.getElementById('popupformWait');
  const popupFormRequestClientStock = document.getElementById('popupformrequestClientStock');
  const popupFormReceiptView = document.querySelectorAll('#popupFormReceiptView');
  const popupFormInvoiceView = document.querySelectorAll('#popupFormInvoiceView');
  const popupFormTransactionView = document.querySelectorAll('#popupFormTransactionView');
  const popupFormUserTransactionView = document.querySelectorAll('#popupFormUserTransactionView');
  const popupFormDebtView = document.querySelectorAll('#popupFormDebtView');

  const deleteButton = document.querySelector('.delete-button');
  const overlay_delete = document.getElementById('overlay_delete');
  const deleteFormContainer = document.querySelector('.delete-form-container');
  const productCloseButtonsedit = document.querySelectorAll('#productCloseButtonedit');
  const closeAddClientStockBtn = document.getElementById('close-btnnnn')
  const reduceClientStockCloseBtn = document.querySelector('#close-buttonReduce');
  const waitCloseBtn = document.querySelector('#close-buttonWait');
  const receiptViewCloseBtns = document.querySelectorAll('#receiptViewCloseBtn');
  const invoiceViewCloseBtns = document.querySelectorAll('#invoiceViewCloseBtn');
  const transactionViewCloseBtns = document.querySelectorAll('#transactionViewCloseBtn');
  const transactionViewUserCloseBtns = document.querySelectorAll('#transactionViewUserCloseBtn');
  const debtViewUserCloseBtns = document.querySelectorAll('#debtViewUserCloseBtn');

  const editButton = document.getElementById('editButton');
  const overlay_edit = document.getElementById('overlayedit');
  const editFormContainer = document.getElementById('popupFormEdit');

  viewInvoiceBtns.forEach((viewInvoiceBtn, index) => {
    viewInvoiceBtn.addEventListener('click', function() {
        overlay_viewInvoices[index].style.display = 'block';
        popupFormInvoiceView[index].style.display = 'block';
    });
});

invoiceViewCloseBtns.forEach((invoiceViewCloseBtn, index) => {
  invoiceViewCloseBtn.addEventListener('click', function() {
    overlay_viewInvoices[index].style.display = 'none';
    popupFormInvoiceView[index].style.display = 'none';
  });
});

viewTransactionBtns.forEach((viewTransactionBtn, index) => {
  viewTransactionBtn.addEventListener('click', function() {
      overlay_viewTransaction[index].style.display = 'block';
      popupFormTransactionView[index].style.display = 'block';
  });
});

transactionViewCloseBtns.forEach((transactionViewCloseBtn, index) => {
  transactionViewCloseBtn.addEventListener('click', function() {
    overlay_viewTransaction[index].style.display = 'none';
    popupFormTransactionView[index].style.display = 'none';
});
});


viewUserTransactionBtns.forEach((viewUserTransactionBtn, index) => {
  viewUserTransactionBtn.addEventListener('click', function() {
      overlay_viewUserTransaction[index].style.display = 'block';
      popupFormUserTransactionView[index].style.display = 'block';
  });
});

transactionViewUserCloseBtns.forEach((transactionViewUserCloseBtn, index) => {
  transactionViewUserCloseBtn.addEventListener('click', function() {
    overlay_viewUserTransaction[index].style.display = 'none';
    popupFormUserTransactionView[index].style.display = 'none';
});
});

viewDebtBtns.forEach((viewDebtBtn, index) => {
  viewDebtBtn.addEventListener('click', function() {
      overlay_viewDebt[index].style.display = 'block';
      popupFormDebtView[index].style.display = 'block';
  });
});

debtViewUserCloseBtns.forEach((debtViewUserCloseBtn, index) => {
  debtViewUserCloseBtn.addEventListener('click', function() {
    overlay_viewDebt[index].style.display = 'none';
    popupFormDebtView[index].style.display = 'none';
});
});

viewReceiptBtns.forEach((viewReceiptBtn, index) => {
    viewReceiptBtn.addEventListener('click', function() {
        overlay_viewReceipts[index].style.display = 'block';
        popupFormReceiptView[index].style.display = 'block';
    });
});

receiptViewCloseBtns.forEach((receiptViewCloseBtn, index) => {
  receiptViewCloseBtn.addEventListener('click', function() {
    overlay_viewReceipts[index].style.display = 'none';
    popupFormReceiptView[index].style.display = 'none';
  });
});

// Selecting request button and close button
const requestButton = document.getElementById('requestClientStockkk');
const requestClientStockCloseBtn = document.getElementById('close-buttonRequest');

// Adding event listener to request button
// if (requestButton) {
//     requestButton.addEventListener('click', function() {
//         console.log("request button clicked");
//         overlay_requestClientStock.style.display = 'block';
//         popupFormRequestClientStock.style.display = 'block';
//     });
// }

// // Adding event listener to close button
// if (requestClientStockCloseBtn) {
//     requestClientStockCloseBtn.addEventListener('click', function(){
//         console.log("close button clicked");
//         overlay_requestClientStock.style.display = 'none';
//         popupFormRequestClientStock.style.display = 'none';
//     });
// }
// });
if (requestButton) {
  requestButton.addEventListener('click', function() {
    overlay_requestClientStock.style.display = 'block';
    popupFormRequestClientStock.style.display = 'block';
  });
}

if (requestClientStockCloseBtn) {
  requestClientStockCloseBtn.addEventListener('click', function(){
    overlay_requestClientStock.style.visibility = 'hidden';
    overlay_requestClientStock.style.opacity = '0';
    popupFormRequestClientStock.style.visibility = 'hidden';
    popupFormRequestClientStock.style.opacity = '0';
  });
}

  if (reduceButton) {
    reduceButton.addEventListener('click', function() {
      overlay_recudeClientStock.style.display = 'block';
      popupFormReduceClientStock.style.display = 'block';
    });
  }
  
  if (reduceClientStockCloseBtn) {
    reduceClientStockCloseBtn.addEventListener('click', function(event){
      event.preventDefault();
      overlay_recudeClientStock.style.visibility = 'hidden';
      overlay_recudeClientStock.style.opacity = '0';
      popupFormReduceClientStock.style.visibility = 'hidden';
      popupFormReduceClientStock.style.opacity = '0';
    });
  }

  if (waitViewButton) {
    waitViewButton.addEventListener('click', function() {
      overlay_Wait.style.display = 'block';
      popupFormWait.style.display = 'block';
    });
  }

  if (waitCloseBtn) {
    waitCloseBtn.addEventListener('click', function(event){
      event.preventDefault();
      overlay_Wait.style.visibility = 'hidden';
      overlay_Wait.style.opacity = '0';
      popupFormWait.style.visibility = 'hidden';
      popupFormWait.style.opacity = '0';
    });
  }

  
  if (closeAddClientStockBtn) {
    closeAddClientStockBtn.addEventListener('click', function(event){
      event.preventDefault();
      overlay_addClientStock.style.visibility = 'hidden';
      overlay_addClientStock.style.opacity = '0';
      popupFormAddClientStock.style.visibility = 'hidden';
      popupFormAddClientStock.style.opacity = '0';
    })
  }

  if (addClientStockBtn) {
    addClientStockBtn.addEventListener('click', function() {
      overlay_addClientStock.style.display = 'block';
      popupFormAddClientStock.style.display = 'block';
    })
  }


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
