from flask import Flask, abort, render_template, redirect, url_for, flash, request
from flask_bootstrap import Bootstrap5
from forms import CreateUserForm, CreateProductForm, EditProductForm, ChangePasswordForm, ResetItForm, ResetPasswordForm, AddAdminStockQuantity, LoginForm, ManageUserForm
# import smtplib
# import os
import requests

HEADERS={
    "Authorization": "Bearer v2.local.jd2hnYhXcDWLwNEebfpBRcp2Pqx983VVbMjwnT7sSRcw_Duw1JDiJaip1m7POgdyx2NPZsTXgTUcKIrkjZo2f102MiNb1xngeI1AB1Kb8O4I0ptJgm5nObdFwVG2pdASSMts8qoQWIfYTnMrPK1zU725BCC171pjy7NdKhHke7609BSBjbxxJKOyNocdnrk6LV4eNxYdAxbHLZatuk7kTEeiTDXqT-TjVmmUA5uTWGW4GZKDs-7IUxtcujjhVR2pfzIEoUo7muGlRGwQ7-C4I7o.bnVsbA"
}

BASE_URL="http://0.0.0.0:8080"


# def send_response(user_name, user_email, user_phone_number, user_message):
#     domain = os.environ.get('DOMAIN')
#     api_key = os.environ.get('APIKEY')
#     mailgun_url = f"https://api.mailgun.net/v3/{domain}/messages"
#     response = requests.post(
#         mailgun_url, 
#         auth=("api", api_key), 
#         data={
#             "from": f"cliff <mailgun@{domain}>", 
#             "to": ["clifftest33@gmail.com"], 
#             "subject": "User Feedback", 
#             "text": f"{user_name} of phone number {user_phone_number} and email {user_email} reached out\n\n{user_message}"
#             }
#         )
   
app = Flask(__name__)
app.config['SECRET_KEY'] = "32e234353t4rffbfbfgxx"
Bootstrap5(app)

@app.route('/create_user"', methods=['GET', 'POST'])
def create_user():
    form = CreateUserForm()
    if form.validate_on_submit():
        createUserUrl = f"{BASE_URL}/users/admin/add"
        createUserRequest = {
            "username":form.username.data,
            "password":"x",
            "email":form.email.data,
            "phone_number":form.phoneNumber.data,
            "address":form.address.data,
            "role":"client",
        }
        rsp = requests.post(url=createUserUrl, json=createUserRequest, headers=HEADERS)
        if rsp.status_code == 200:
            flash("Login successful")
            return redirect(url_for('list_users'))
        else:
            return render_template('failed.html', error_code=rsp.status_code)
    return render_template("create_user.html", form=form)

@app.route('/create_product"', methods=['GET', 'POST'])
def create_product():
    form = CreateProductForm()
    if form.validate_on_submit():
        createProductUrl = f"{BASE_URL}/products/admin/add"
        createProductRequest = {
            "product_name": form.productName.data,
            "unit_price": int(form.unitPrice.data),
            "packsize": form.packsize.data
        }
        rsp = requests.post(url=createProductUrl, json=createProductRequest, headers=HEADERS)
        if rsp.status_code == 200:
            flash("Created Successfully")
            return redirect(url_for('list_products'))
        else:
            return render_template('failed.html', error_code=rsp.status_code)
    return render_template("create_product.html", form=form)

@app.route('/delete_product/<int:id>')
def delete_product(id):
    deleteProductUrl = f"{BASE_URL}/products/admin/delete/{id}"
    rsp = requests.delete(url=deleteProductUrl, headers=HEADERS)
    if rsp.status_code == 200:
        flash("Product Deleted Successfully")
        return redirect(url_for('get_product'))
    else:
        return render_template('failed.html', error_code=rsp.status_code) 

@app.route('/delete_user/<int:id>')
def delete_user(id):
    deleteUserUrl = f"http://0.0.0.0:8080/users/admin/{id}"
    rsp = requests.delete(url=deleteUserUrl, headers=HEADERS)
    if rsp.status_code == 200:
        flash("User Deleted Successfully")
        return redirect(url_for('get_user'))
    else:
        return render_template('failed.html', error_code=rsp.status_code)       
    
@app.route('/edit_product/<int:id>', methods=['GET', 'POST'])
def edit_product(id):
    getProductUrl = f"{BASE_URL}/products/{id}"
    product = requests.get(url=getProductUrl, headers=HEADERS)
    data = product.json()
    print(product.text)
    form = EditProductForm(
        productName=data['product_name'],
        unitPrice=int(data['unit_price']),
        packSize=data['packsize'],
    )
    if form.validate_on_submit():
        editProductUrl = f"{BASE_URL}/products/admin/edit/{id}"
        editProductRequest = {
            "product_name": form.productName.data,
            "unit_price": form.unitPrice.data,
            "packsize": form.packSize.data
        }
        rsp = requests.put(url=editProductUrl, json=editProductRequest, headers=HEADERS)
        if rsp.status_code == 200:
            flash("Users Details Changed Successfully")
            return redirect(url_for('get_product'))
        else:
            return render_template('failed.html', error_code=rsp.status_code)          
    return render_template("edit_product.html", form=form)

@app.route('/change_password/<int:id>', methods=['GET', 'POST'])
def change_password(id):
    form = ChangePasswordForm()
    if form.validate_on_submit():
        changePasswordUrl = f"{BASE_URL}/users/{id}/edit"
        changePasswordRequest = {
            "old_password": form.oldPassword.data,
            "new_password": form.newPassword.data,
            "role": "client"
        }
        rsp = requests.put(url=changePasswordUrl, json=changePasswordRequest, headers=HEADERS)
        if rsp.status_code == 200:
            flash("Users Details Changed Successfully")
            return redirect(url_for('get_user'))
        else:
            return render_template('failed.html', error_code=rsp.status_code)          
    return render_template("change_password.html", form=form)

@app.route('/manage_user/<int:id>', methods=['GET', 'POST'])
def manage_user(id):
    user = requests.get(url=f"{BASE_URL}/users/{id}", headers=HEADERS)
    data = user.json()
    form = ManageUserForm(
        email=data['email'],
        phoneNumber=data['phone_number'],
        address=data['address'],
        username=data['username']
    )
    if form.validate_on_submit():
        manageUserUrl = f"{BASE_URL}/users/admin/manage/{id}"
        changePasswordRequest = {
            "email": form.email.data,
            "phone_number": form.phoneNumber.data,
            "address": form.address.data,
            "username": form.username.data,
            "role": "admin"
        }
        rsp = requests.put(url=manageUserUrl, json=changePasswordRequest, headers=HEADERS)

        if rsp.status_code == 200:
            flash("Users Details Changed Successfully")
            return redirect(url_for('get_user'))
        else:
            return render_template('failed.html', error_code=rsp.status_code)           
    return render_template("manage_user.html", form=form)

@app.route("/reset", methods=['GET', 'POST'])
def reset():
    form = ResetPasswordForm()
    if form.validate_on_submit():
        resetPasswordUrl = f"{BASE_URL}/reset"
        resetPasswordRequest = {"email": form.email.data}
        rsp = requests.post(url=resetPasswordUrl, json=resetPasswordRequest, headers=HEADERS)
        if rsp.status_code == 200:
            return redirect(url_for('login'))
        else:
            return render_template('failed.html', error_code=rsp.status_code)     
    return render_template("reset_password.html", form=form)

@app.route("/resetit", methods=['GET', 'POST'])
def resetit():
    form = ResetItForm()
    if form.validate_on_submit():
        resetItUrl = f"{BASE_URL}/resetit"
        rsp = requests.post(url=resetItUrl, params={"token": request.args.get('token')}, json={"password": form.password.data}, headers=HEADERS)
        if rsp.status_code == 200:
            return redirect(url_for('dashboard'))
        else:
            return render_template('failed.html', error_code=rsp.status_code)
    return render_template("reset_it.html", form=form)
    
@app.route('/login_user', methods=['GET', 'POST'])
def login():
    form = LoginForm()
    if form.validate_on_submit():
        userLoginRequest = {
            "email": form.email.data,
            "password": form.password.data
        }
        userLoginUrl = f"{BASE_URL}/users/login"
        rsp = requests.get(url=userLoginUrl, json=userLoginRequest, headers=HEADERS)
        if rsp.status_code == 200:
            return redirect(url_for('dashboard'))
        else:
            return render_template('failed.html', error_code=rsp.status_code)
    return render_template("login_user.html", form=form)

@app.route('/get_user/<int:id>')
def get_user(id):
    getUserUri = f"{BASE_URL}/users/{id}"
    rsp = requests.get(url=getUserUri, headers=HEADERS)
    if rsp.status_code == 200:
        return render_template('get_user.html', user_data=rsp.json()) # unmarshal JSON and read data
    else:
        return render_template('failed.html', error_code=rsp.status_code)

@app.route('/get_product/<int:id>')
def get_product(id):
    getProductUri = f"{BASE_URL}/products/{id}"
    rsp = requests.get(url=getProductUri, headers=HEADERS)
    if rsp.status_code == 200:
        return render_template('get_prouct.html', user_data=rsp.json()) # unmarshal JSON and read data
    else:
        return render_template('failed.html', error_code=rsp.status_code)

@app.route('/get_invoice/<string:invoice_number>')
def get_invoice(invoice_number):
    getInvoiceUri = f"{BASE_URL}/invoices/{invoice_number}"
    rsp = requests.get(url=getInvoiceUri, headers=HEADERS)
    if rsp.status_code == 200:
        return render_template('get_invoice.html', user_data=rsp.json()) # unmarshal JSON and read data
    else:
        return render_template('failed.html', error_code=rsp.status_code)
    
@app.route('/get_receipt/<int:id>')
def get_receipt(id):
    getReceiptUri = f"{BASE_URL}/receipts/{id}"
    rsp = requests.get(url=getReceiptUri, headers=HEADERS)
    if rsp.status_code == 200:
        return render_template('get_receipt.html', user_data=rsp.json()) # unmarshal JSON and read data
    else:
        return render_template('failed.html', error_code=rsp.status_code)

@app.route('/')
def dashboard():
    return render_template("index.html", dataToDisplay="Home Page") # Figure what should be on the dashboard and pass its data



@app.route('/list_invoices')
def list_invoices():
    listInvoicesUri = f"{BASE_URL}/invoices/admin"
    rsp = requests.get(url=listInvoicesUri, headers=HEADERS)
    if rsp.status_code != 200:
        flash("Please try again server error")
        return render_template('failed.html', error_code=rsp.status_code)
    return render_template("list_invoices.html", invoices_data=rsp.json())

@app.route('/list_receipts')
def list_receipts():
    listReceiptUrl = f"{BASE_URL}/receipts/admin"
    rsp = requests.get(url=listReceiptUrl, headers=HEADERS)
    if rsp.status_code != 200:
        flash("Please try again server error")
        return render_template('failed.html', error_code=rsp.status_code)
    return render_template("list_receipts.html", receipts_data=rsp.json())

@app.route('/list_users')
def list_users():
    listUsersUrl = f"{BASE_URL}/users/admin"
    rsp = requests.get(url=listUsersUrl, headers=HEADERS)
    if rsp.status_code != 200:
        flash("Please try again server error")
        return render_template('failed.html', error_code=rsp.status_code)
    return render_template("list_users.html", users_data=rsp.json())

@app.route('/list_products')
def list_products():
    listProductsUrl = f"{BASE_URL}/products"
    rsp = requests.get(url=listProductsUrl, headers=HEADERS)
    if rsp.status_code != 200:
        flash("Please try again server error")
        return render_template('failed.html', error_code=rsp.status_code)
    return render_template("list_products.html", products_data=rsp.json())

@app.route('/get_user_invoices/<int:id>')
def get_user_invoices(id):
    getUserInvoiceUrl = f"{BASE_URL}/users/invoices/{id}"
    rsp = requests.get(url=getUserInvoiceUrl, headers=HEADERS)
    if rsp.status_code == 200:
        return render_template("get_user_invoices.html", user_invoiceData=rsp.json())
    else:
        return render_template('failed.html', error_code=rsp.status_code)
    
@app.route('/get_user_receipts/<int:id>')
def get_user_receipts(id):
    getUserReceiptsUrl = f"{BASE_URL}/users/receipts/{id}"
    rsp = requests.get(url=getUserReceiptsUrl, headers=HEADERS)
    if rsp.status_code == 200:
        return render_template("get_user_receipts.html", user_invoiceData=rsp.json())
    else:
        return render_template('failed.html', error_code=rsp.status_code)
    
@app.route('/get_user_products/<int:id>')
def get_user_products(id):
    getUserProductsUrl = f"{BASE_URL}/users/products/{id}"
    rsp = requests.get(url=getUserProductsUrl, headers=HEADERS)
    if rsp.status_code == 200:
        return render_template("get_user_products.html", user_invoiceData=rsp.json())
    else:
        return render_template('failed.html', error_code=rsp.status_code)

@app.route("/users/admin/manage/add/<int:id>", methods=['POST', 'GET'])
def add_admin_stock(id):
    form = AddAdminStockQuantity()
    if form.validate_on_submit():
        addAdminStockUrl = f"{BASE_URL}/users/products/admin/add/{id}"
        q = form.quantity.data
        rsp = requests.post(url=addAdminStockUrl, json={"user_id": 1, "quantity": q}, headers=HEADERS)
        if rsp.status_code == 200:
            return redirect(url_for('get_product'))
        else:
            return render_template('failed.html', error_code=rsp.status_code)
    return render_template("add_admin_stock.html", form=form)

@app.route('/add_client_stock/<int:id>', methods=['POST', 'GET'])
def add_client_stock(id):
    if request.method == 'POST':
        products_id = request.form.get('products_id')
        product_list = [int(num) for num in products_id.split(",")]
        quantities = request.form.get('quantities')
        quantities_list = [int(num) for num in quantities.split(",")]
        print(type(products_id))
        print(products_id)

        data = {
            "products_id": product_list,
            "quantities": quantities_list
        }
        addClientStockUrl = f"{BASE_URL}/users/admin/manage/add/{id}"
        rsp = requests.post(url=addClientStockUrl, json=data, headers=HEADERS)
        if rsp.status_code == 200:
            return redirect(url_for('get_product'))
        else:
            return render_template('failed.html', error_code=rsp.status_code)
    return render_template("add_client_stock.html")

@app.route('/users/products/sell/<int:id>', methods=['POST', 'GET'])
def reduce_client_stock(id):
    if request.method == 'POST':
        products_id = request.form.get('products_id')
        product_list = [int(num) for num in products_id.split(",")]
        quantities = request.form.get('quantities')
        quantities_list = [int(num) for num in quantities.split(",")]

        data = {
            "products_id": product_list,
            "quantities": quantities_list
        }

        url = f"{BASE_URL}/users/products/sell/{id}"
        rsp = requests.post(url, json=data, headers=HEADERS)

        if rsp.status_code == 200:
            return redirect(url_for('get_product'))  # Redirect to some success page
        else:
            return render_template('failed.html', error_code=rsp.status_code)
    return render_template("reduce_client_stock.html")

if __name__ == "__main__":
    app.run(debug=True)
