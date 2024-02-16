from flask import Flask, abort, render_template, redirect, url_for, flash, request, session
from flask_bootstrap import Bootstrap5
from forms import CreateUserForm, CreateProductForm, EditProductForm, ChangePasswordForm, ResetItForm, ResetPasswordForm, AddAdminStockQuantity, LoginForm, ManageUserForm
# import smtplib
# import os
import requests

HEADERS={
    "Authorization": "Bearer "
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
    if request.method  == "POST":
        createUserUrl = f"{BASE_URL}/users/admin/add"
        createUserRequest = {
            "username": request.form.get('username'),
            "password":"x",
            "email": request.form.get('email'),
            "phone_number": request.form.get('phone'),
            "address": request.form.get('address'),
            "role":"client",
        }
        rsp = requests.post(url=createUserUrl, json=createUserRequest, headers={"Authorization": f"Bearer {session['token']}"})
        if rsp.status_code == 200:
            flash("Login successful")
            return redirect(url_for('list_users'))
        else:
            return render_template('failed.html', error_code=rsp.status_code)
    return redirect(url_for("list_users.html"))

@app.route('/create_product"', methods=['GET', 'POST'])
def create_product():
    if request.method == "POST":
        createProductUrl = f"{BASE_URL}/products/admin/add"
        createProductRequest = {
            "product_name": request.form.get('product_name'),
            "unit_price": int(request.form.get('unit_price')),
            "packsize": request.form.get('packsize')
        }
        rsp = requests.post(url=createProductUrl, json=createProductRequest, headers={"Authorization": f"Bearer {session['token']}"})
        if rsp.status_code == 200:
            flash("Created Successfully")
            return redirect(url_for('list_products'))
        else:
            return render_template('failed.html', error_code=rsp.status_code)
    return redirect(url_for("list_users.html"))

@app.route('/delete_product/<int:id>', methods=['POST'])
def delete_product(id):
    deleteProductUrl = f"{BASE_URL}/products/admin/delete/{id}"
    rsp = requests.delete(url=deleteProductUrl, headers={"Authorization": f"Bearer {session['token']}"})
    if rsp.status_code == 200:
        flash("Product Deleted Successfully")
        return redirect(url_for('list_products'))
    else:
        return render_template('failed.html', error_code=rsp.status_code) 

@app.route('/delete_user/<int:id>', methods=['POST'])
def delete_user(id):
    if request.method == "POST":
        getUserUri = f"{BASE_URL}/users/{id}"
        rsp = requests.get(url=getUserUri, headers={"Authorization": f"Bearer {session['token']}"})
        user = rsp.json()
        print(user)
        if user['username'].lower() == request.form.get('delete-username').lower():
            deleteUserUrl = f"http://0.0.0.0:8080/users/admin/{id}"
            rsp = requests.delete(url=deleteUserUrl, headers={"Authorization": f"Bearer {session['token']}"})
            if rsp.status_code == 200:
                flash("User Deleted Successfully")
                return redirect(url_for('list_users'))
            else:
                return render_template('failed.html', error_code=rsp.status_code) 
        else:
            return render_template('failed.html', error_code={"failed":"incorrect user name"})     
    
@app.route('/edit_product/<int:id>', methods=['GET', 'POST'])
def edit_product(id):
    getProductUrl = f"{BASE_URL}/products/{id}"
    product = requests.get(url=getProductUrl, headers={"Authorization": f"Bearer {session['token']}"})
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
        rsp = requests.put(url=editProductUrl, json=editProductRequest, headers={"Authorization": f"Bearer {session['token']}"})
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
        rsp = requests.put(url=changePasswordUrl, json=changePasswordRequest, headers={"Authorization": f"Bearer {session['token']}"})
        if rsp.status_code == 200:
            flash("Users Details Changed Successfully")
            return redirect(url_for('get_user'))
        else:
            return render_template('failed.html', error_code=rsp.status_code)          
    return render_template("change_password.html", form=form)

@app.route('/manage_user/<int:id>', methods=['GET', 'POST'])
def manage_user(id):
    user = requests.get(url=f"{BASE_URL}/users/{id}", headers={"Authorization": f"Bearer {session['token']}"})
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
        rsp = requests.put(url=manageUserUrl, json=changePasswordRequest, headers={"Authorization": f"Bearer {session['token']}"})

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
        rsp = requests.post(url=resetPasswordUrl, json=resetPasswordRequest, headers={"Authorization": f"Bearer {session['token']}"})
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
        rsp = requests.post(url=resetItUrl, params={"token": request.args.get('token')}, json={"password": form.password.data}, headers={"Authorization": f"Bearer {session['token']}"})
        if rsp.status_code == 200:
            return redirect(url_for('dashboard'))
        else:
            return render_template('failed.html', error_code=rsp.status_code)
    return render_template("reset_it.html", form=form)
    
@app.route('/login_user', methods=['GET', 'POST'])
def login():
    if request.method == "POST":
        userLoginRequest = {
            "email": request.form['email'],
            "password": request.form['pass']
        }
        userLoginUrl = f"{BASE_URL}/users/login"
        rsp = requests.get(url=userLoginUrl, json=userLoginRequest, headers={"Authorization": f"Bearer {session['token']}"})
        user_response = rsp.json()
        session['token'] = user_response['access_token']
        session['user_id'] = user_response['user']['id']
        if rsp.status_code == 200:
            return redirect(url_for('dashboard'))
        elif rsp.status_code == 400:
            render_template('failed.html', error_code=rsp.status_code)
        elif rsp.status_code == 401:
            flash("Incorrect Password")
        elif rsp.status_code == 404:
            flash("No user with this email found")
        else:
            return render_template('failed.html', error_code=rsp.status_code)
    return render_template("login.html")

@app.route('/get_user/<int:id>')
def get_user(id):
    getUserUri = f"{BASE_URL}/users/{id}"
    rsp = requests.get(url=getUserUri, headers={"Authorization": f"Bearer {session['token']}"})
    if rsp.status_code == 200:
        return render_template('user.html', user=rsp.json(), user_id=session['user_id'], ct="user") # unmarshal JSON and read data
    else:
        return render_template('failed.html', error_code=rsp.status_code)

@app.route('/get_product/<int:id>')
def get_product(id):
    getProductUri = f"{BASE_URL}/products/{id}"
    rsp = requests.get(url=getProductUri, headers={"Authorization": f"Bearer {session['token']}"})
    if rsp.status_code == 200:
        return render_template('get_product.html', user_data=rsp.json()) # unmarshal JSON and read data
    else:
        return render_template('failed.html', error_code=rsp.status_code)

@app.route('/get_invoice/<int:id>')
def get_invoice(id):
    getInvoiceUri = f"{BASE_URL}/invoices/{id}"
    rsp = requests.get(url=getInvoiceUri, headers={"Authorization": f"Bearer {session['token']}"})
    if rsp.status_code == 200:
        return render_template('get_invoice.html', user_data=rsp.json()) # unmarshal JSON and read data
    else:
        return render_template('failed.html', error_code=rsp.status_code)
    
@app.route('/get_receipt/<int:id>')
def get_receipt(id):
    getReceiptUri = f"{BASE_URL}/receipts/{id}"
    rsp = requests.get(url=getReceiptUri, headers={"Authorization": f"Bearer {session['token']}"})
    if rsp.status_code == 200:
        return render_template('get_receipt.html', user_data=rsp.json()) # unmarshal JSON and read data
    else:
        return render_template('failed.html', error_code=rsp.status_code)

@app.route('/')
def dashboard():
    session['user_id'] = 1
    return render_template("index.html", user_id=session['user_id'])

@app.route('/list_invoices')
def list_invoices():
    listInvoicesUri = f"{BASE_URL}/invoices/admin"
    rsp = requests.get(url=listInvoicesUri, headers={"Authorization": f"Bearer {session['token']}"})
    if rsp.status_code != 200:
        flash("Please try again server error")
        return render_template('failed.html', error_code=rsp.status_code)
    return render_template("list.html", data_sent=rsp.json(), ct="invoices", user_id=session['user_id'])

@app.route('/list_receipts')
def list_receipts():
    listReceiptUrl = f"{BASE_URL}/receipts/admin"
    rsp = requests.get(url=listReceiptUrl, headers={"Authorization": f"Bearer {session['token']}"})
    if rsp.status_code != 200:
        flash("Please try again server error")
        return render_template('failed.html', error_code=rsp.status_code)
    return render_template("list.html", data_sent=rsp.json(), ct="receipts", user_id=session['user_id'])

@app.route('/list_users', methods=['GET', 'POST'])
def list_users():
    listUsersUrl = f"{BASE_URL}/users/admin"
    rsp = requests.get(url=listUsersUrl, headers={"Authorization": f"Bearer {session['token']}"})
    if rsp.status_code != 200:
        flash("Please try again server error")
        return render_template('failed.html', error_code=rsp.status_code)
    return render_template("list.html", data_sent=rsp.json(), ct="users", user_id=session['user_id'])

@app.route('/search_users', methods=['GET', 'POST'])
def search_users():
    if request.method == 'POST':
        query = request.form.get('search')
        listUsersUrl = f"{BASE_URL}/search/users"
    rsp = requests.get(url=listUsersUrl, json={"search_word": query}, headers={"Authorization": f"Bearer {session['token']}"})
    if rsp.status_code != 200:
        flash("Please try again server error")
        return render_template('failed.html', error_code=rsp.status_code)
    return render_template("list.html", data_sent=rsp.json(), ct="users", user_id=session['user_id'])

@app.route('/search_products', methods=['GET', 'POST'])
def search_products():
    if request.method == 'POST':
        query = request.form.get('search')
        listUsersUrl = f"{BASE_URL}/search/products"
    rsp = requests.get(url=listUsersUrl, json={"search_word": query}, headers={"Authorization": f"Bearer {session['token']}"})
    if rsp.status_code != 200:
        flash("Please try again server error")
        return render_template('failed.html', error_code=rsp.status_code)
    return render_template("list.html", data_sent=rsp.json(), ct="products", user_id=session['user_id'])

@app.route('/list_products')
def list_products():
    listProductsUrl = f"{BASE_URL}/products"
    rsp = requests.get(url=listProductsUrl, headers={"Authorization": f"Bearer {session['token']}"})
    if rsp.status_code != 200:
        flash("Please try again server error")
        return render_template('failed.html', error_code=rsp.status_code)
    return render_template("list.html", data_sent=rsp.json(), ct="products", user_id=session['user_id'])

@app.route('/get_user_invoices/<int:id>')
def get_user_invoices(id):
    user = requests.get(url=f"{BASE_URL}/users/{id}", headers={"Authorization": f"Bearer {session['token']}"})
    getUserInvoiceUrl = f"{BASE_URL}/users/invoices/{id}"
    rsp = requests.get(url=getUserInvoiceUrl, headers={"Authorization": f"Bearer {session['token']}"})
    if rsp.status_code == 200:
        return render_template("user.html", invoice=rsp.json(), user=user.json(), user_id=session['user_id'], ct='invoice')
    else:
        return render_template('failed.html', error_code=rsp.status_code)
    
@app.route('/get_user_receipts/<int:id>')
def get_user_receipts(id):
    user = requests.get(url=f"{BASE_URL}/users/{id}", headers={"Authorization": f"Bearer {session['token']}"})
    getUserReceiptsUrl = f"{BASE_URL}/users/receipts/{id}"
    rsp = requests.get(url=getUserReceiptsUrl, headers={"Authorization": f"Bearer {session['token']}"})
    if rsp.status_code == 200:
        return render_template("user.html", receipt=rsp.json(), user=user.json(), user_id=session['user_id'], ct='receipt')
    else:
        return render_template('failed.html', error_code=rsp.status_code)
    
@app.route('/get_user_products/<int:id>')
def get_user_products(id):
    getUserProductsUrl = f"{BASE_URL}/users/products/{id}"
    rsp = requests.get(url=getUserProductsUrl, headers={"Authorization": f"Bearer {session['token']}"})
    if rsp.status_code == 200:
        return render_template("get_user_products.html", user_invoiceData=rsp.json())
    else:
        return render_template('failed.html', error_code=rsp.status_code)

@app.route("/users/admin/manage/add", methods=['POST', 'GET'])
def add_admin_stock():
    if request.method == 'POST':
        id = request.form.get('productID')
        print(id)
        addAdminStockUrl = f"{BASE_URL}/users/products/admin/add/{id}"
        q = request.form.get("quantity")
        rsp = requests.post(url=addAdminStockUrl, json={"user_id": 1, "quantity": int(q)}, headers={"Authorization": f"Bearer {session['token']}"})
        if rsp.status_code == 200:
            return redirect(url_for('get_user', id=1))
        else:
            return render_template('failed.html', error_code=rsp.status_code)

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
        rsp = requests.post(url=addClientStockUrl, json=data, headers={"Authorization": f"Bearer {session['token']}"})
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
        rsp = requests.post(url, json=data, headers={"Authorization": f"Bearer {session['token']}"})

        if rsp.status_code == 200:
            return redirect(url_for('get_product'))  # Redirect to some success page
        else:
            return render_template('failed.html', error_code=rsp.status_code)
    return render_template("reduce_client_stock.html")

if __name__ == "__main__":
    app.run(debug=True)
