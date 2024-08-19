const adminDataDiv = document.getElementById('adminDataDiv');
    
    const adminProducts = JSON.parse(adminDataDiv.dataset.adminData);
    
    let nextBtn = document.getElementById('nextBtn');
    let holder = document.querySelector(".container-form-holder");
    let supplierForm = document.getElementById('supplierForm');
    supplierForm.addEventListener('submit', function(event) {
        event.preventDefault();
    });

    nextBtn.addEventListener('click', function() {
        // supplierForm.preventDefault();
        const formData = new FormData(supplierForm);
        let values = {};
    
        for (const [key, value] of formData.entries()) {
            values[key] = value;
        }
    
        localStorage.setItem('supplier', JSON.stringify(values));
        holder.innerHTML = '';
        holder.innerHTML = `
        <p>Purchase Order Products</p>
        <button type="button" class="btn btn-sm btn-primary" onclick="addOrderProduct()">Add</button>
        <form id="productForm" autocomplete="off">
            <div id="productList" class="d-flex flex-column form-group form-floating m-1 p-1" style="overflow: scroll; row-gap: .25rem;">
                <div class="d-flex">
                    <input type="text" class="productName_0" name="productName_0" id="productName_0" placeholder="Product Name" style="width: 12rem; margin-right: 2rem" required>
                    <input type="number" class="quantity_0" min="1" name="quantities_0" id="quantities_0" placeholder="Quantity" style="width: 5.2rem; margin-right: 1rem" required>
                    <input type="number" class="price_0"  min="1" name="price_0" id="price_0" placeholder="Price" style="width: 5rem;" required>
                </div>
            </div>
            <button type="" id="finalSubmit" class="btn btn-primary">Submit</button>
        </form>
        `;
    
        let dataForm = document.getElementById("productForm")
        dataForm.addEventListener('submit', function(event) {
            event.preventDefault();
        });
    
        let submitBtn = document.getElementById("finalSubmit");
        submitBtn.addEventListener('click', function() {
            const formData = new FormData(dataForm);
            let values = {};
            
            for (const [key, value] of formData.entries()) {
                values[key] = value;
            }
            supplierData = JSON.parse(localStorage.getItem('supplier'));

            let sendData = {
                "supplier_name": supplierData.supplier,
                "po_box":  supplierData.poBox,
                "address": supplierData.address,
                "data": []
            }

            let empty = false;
            for (let i = 0; i <= productCount; i++) {
                if (Number(values[`quantities_${i}`]) === 0 || Number(values[`price_${i}`]) === 0) {
                    empty = true;
                    break
                } else {
                    sendData.data.push({
                        "product_name": values[`productName_${i}`],
                        "quantity":     Number(values[`quantities_${i}`]),
                        "unit_price":   Number(values[`price_${i}`])
                    })

                }
            }

            if (empty) {
                alert("Please fill all the fields")
            } else {
                // sendDataFunc(sendData)
                sendDataFuncDirect(sendData)
            }
            
            // console.log(sendData)
        })
    })

function sendDataFuncDirect(sendData) {
    const form = document.createElement('form');
    form.method = 'POST';
    form.action = '/download/purchase-order';
    form.target = '_blank';

    const input = document.createElement('input');
    input.type = 'hidden';
    input.name = 'data';
    input.value = JSON.stringify(sendData);
    form.appendChild(input);

    document.body.appendChild(form);
    form.submit();
    document.body.removeChild(form);

    setTimeout(function() {
        window.location.href = "/get_user/1";
    }, 1000);
}

function sendDataFunc(sendData) {
    fetch("/download/purchase-order", {
        "method": "POST",
        "headers": {"Content-Type": "application/json"},
        "body": JSON.stringify(sendData),
    })
    .then(response => response.blob())
    .then(blob => {
        // Create a link element to download the file
        console.log(blob)
        const url = window.URL.createObjectURL(blob);
        const a = document.createElement('a');
        a.href = url;
        a.download = 'purchase_order.pdf';
        document.body.appendChild(a);
        a.click();
        a.remove();
        
        window.location.href ="/get_user/1"
    })
    .catch(error => console.error('Error:', error));
}

let productCount = 0; 

function addOrderProduct() {
    productCount++;
        var productName = document.createElement('input');
        productName.style.width = '12rem';
        productName.type = 'text';
        productName.style.marginRight = '2rem';
        productName.className = `productName_${productCount}`;
        productName.name = `productName_${productCount}`;
        productName.id = `productName_${productCount}`;
        productName.placeholder = 'Price';
        productName.required = true;
    
        var newQuantity = document.createElement('input');
        newQuantity.style.width = '5.2rem';
        newQuantity.style.marginRight = '1rem';
        newQuantity.type = 'number';
        newQuantity.className = `quantity_${productCount}`;
        newQuantity.name = `quantities_${productCount}`;
        newQuantity.id = `quantities_${productCount}`;
        newQuantity.placeholder = 'Quantity';
        newQuantity.required = true;
        newQuantity.min = '1';

        var newPrice = document.createElement('input');
        newPrice.style.width = '5rem';
        newPrice.type = 'number';
        newPrice.className = `price_${productCount}`;
        newPrice.name = `price_${productCount}`;
        newPrice.id = `price_${productCount}`;
        newPrice.placeholder = 'Price';
        newPrice.required = true;
        newPrice.min = '1';
    
        var newProductDiv = document.createElement('div');
        newProductDiv.className = 'd-flex';
        newProductDiv.appendChild(productName);
        newProductDiv.appendChild(newQuantity);
        newProductDiv.appendChild(newPrice);
    
        document.getElementById('productList').appendChild(newProductDiv);
}