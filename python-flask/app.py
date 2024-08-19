from flask import Flask, abort, render_template, redirect, url_for, flash, request, session, send_file
from flask_bootstrap import Bootstrap5
from forms import ChangePasswordForm
import requests
import base64
import json
from io import BytesIO
from requests.exceptions import ConnectionError
from werkzeug.exceptions import InternalServerError
from collections import OrderedDict

HEADERS={
    "Authorization": "Bearer "
}

# BASE_URL="http://backend:8080" # When Testing
# BASE_URL = "http://inventory-system-api-1:8080" When using Docker Compose
BASE_URL = "http://secretive-window.railway.internal:8080"  #  Production
   
app = Flask(__name__)
app.config['SECRET_KEY'] = "32e234353t4rffbfbfgxx"
app.config['SESSION_PERMANENT'] = True
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
            flash("User Created Succefully", "success")
            return redirect(url_for('list_users'))
        elif rsp.status_code == 401:
            flash("Please login", "error")
            return redirect(url_for('login'))
        elif rsp.status_code == 500:
            if rsp.json()["error"] == "ERROR: duplicate key value violates unique constraint \"users_email_key\" (SQLSTATE 23505)":
                flash("email/phone_number already exists", "error")
                return redirect(url_for('list_users'))
            return render_template('failed.html', error_code=rsp.status_code, error=rsp.json()['error'])
        else:
            return render_template('failed.html', error_code=rsp.status_code, error=rsp.json()['error'])
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
            flash("Product Created Successfully", "success")
            return redirect(url_for('list_products'))
        if rsp.status_code == 409:
            flash("Product Already Exists", "error")
            return redirect(url_for('list_products'))
        elif rsp.status_code == 401:
            flash("Please login", "error")
            return redirect(url_for('login'))
        elif rsp.status_code == 500:
            if rsp.json()["error"] == "ERROR: duplicate key value violates unique constraint \"products_product_name_key\" (SQLSTATE 23505)":
                flash("Product Already Exists", "error")
                return redirect(url_for('list_products'))
            return render_template('failed.html', error_code=rsp.status_code, error=rsp.json()['error'])
        else:
            return render_template('failed.html', error_code=rsp.status_code, error=rsp.json()['error'])
    return redirect(url_for("list_products.html"))

@app.route('/delete_product/<int:id>', methods=['POST'])
def delete_product(id):
    deleteProductUrl = f"{BASE_URL}/products/admin/delete/{id}"
    rsp = requests.delete(url=deleteProductUrl, headers={"Authorization": f"Bearer {session['token']}"})
    if rsp.status_code == 200:
        flash("Product Deleted Successfully", "success")
        return redirect(url_for('list_products'))
    elif rsp.status_code == 401:
        flash("Please login", "error")
        return redirect(url_for('login'))
    else:
        return render_template('failed.html', error_code=rsp.status_code, error=rsp.json()['error'])

@app.route('/delete_user/<int:id>', methods=['POST'])
def delete_user(id):
    if request.method == "POST":
        getUserUri = f"{BASE_URL}/users/{id}"
        rsp = requests.get(url=getUserUri, headers={"Authorization": f"Bearer {session['token']}"})
        user = rsp.json()
        print(user)
        if user['username'].lower() == request.form.get('delete-username').lower():
            deleteUserUrl = f"{BASE_URL}/users/admin/{id}"
            rsp = requests.delete(url=deleteUserUrl, headers={"Authorization": f"Bearer {session['token']}"})
            if rsp.status_code == 200:
                flash("User Deleted Successfully", "success")
                return redirect(url_for('list_users'))
            elif rsp.status_code == 401:
                flash("Please login", "error")
                return redirect(url_for('login'))
            else:
                return render_template('failed.html', error_code=rsp.status_code, error=rsp.json()['error'])
        else:
            flash("Username didnt match", "error")
            return redirect(url_for('list_users'))
    
@app.route('/edit_product/<int:id>', methods=['POST'])
def edit_product(id):
    if request.method == "POST":
        editProductUrl = f"{BASE_URL}/products/admin/edit/{id}"
        editProductRequest = {
            "product_name": request.form.get("product_name"),
            "unit_price": int(request.form.get("unit_price")),
            "packsize": request.form.get("packsize")
        }
        print(editProductRequest)
        rsp = requests.put(url=editProductUrl, json=editProductRequest, headers={"Authorization": f"Bearer {session['token']}"})
        if rsp.status_code == 200:
            flash("Products Details Changed Successfully", "success")
            return redirect(url_for('list_products'))
        elif rsp.status_code == 409:
            flash("Products Already Exists", "error")
            return redirect(url_for('login'))
        elif rsp.status_code == 401:
            flash("Please login", "error")
            return redirect(url_for('login'))
        else:
            return render_template('failed.html', error_code=rsp.status_code, error=rsp.json()['error'])         
    return render_template("edit_product.html")

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
            flash("Password Details Changed Successfully", "success")
            return redirect(url_for('get_user', id=id))
        elif rsp.status_code == 401:
            flash("Please login", "error")
            return redirect(url_for('login'))
        else:
            return render_template('failed.html', error_code=rsp.status_code, error=rsp.json()['error'])
    return render_template("change_password.html", form=form)

@app.route('/manage_user/<int:id>', methods=['GET', 'POST'])
def manage_user(id):
    if request.method == "POST":
        manageUserUrl = f"{BASE_URL}/users/admin/manage/{id}"
        changePasswordRequest = {
            "email": request.form.get("email"),
            "phone_number": request.form.get("phone"),
            "address": request.form.get("address"),
            "username": request.form.get("username"),
            "role": "admin"
        }
        print(changePasswordRequest)
        rsp = requests.put(url=manageUserUrl, json=changePasswordRequest, headers={"Authorization": f"Bearer {session['token']}"})

        if rsp.status_code == 200:
            flash("Users Details Changed Successfully", "success")
            return redirect(url_for('get_user', id=id))
        elif rsp.status_code == 401:
            flash("Please login", "error")
            return redirect(url_for('login'))
        else:
            return render_template('failed.html', error_code=rsp.status_code, error=rsp.json()['error'])          
    return redirect(url_for("get_user.html", id))

@app.route("/reset", methods=['GET', 'POST'])
def reset():
    if request.method == "POST":
        resetPasswordUrl = f"{BASE_URL}/reset"
        resetPasswordRequest = {"email": request.form.get('email')}
        rsp = requests.post(url=resetPasswordUrl, json=resetPasswordRequest)
        if rsp.status_code == 200:
            flash("Reset email sent", "success")
            return redirect(url_for('login'))
        elif rsp.status_code == 500:
            flash("No User Found With Email Provided", "error")
            return redirect(url_for('reset'))
        else:
            return render_template('failed.html', error_code=rsp.status_code, error=rsp.json()['error'])    
    return render_template("forgot_password.html", reset=True)

@app.route("/resetit", methods=['GET', 'POST'])
def resetit():
    token = request.args.get('token')
    if request.method == "POST":
        token = request.form.get('token')
        password = request.form.get('pass')
        confimPass = request.form.get('Confirmpass')
        if password != confimPass:
            flash("Password Don't Match", "error")
            return redirect('reset')
        resetItUrl = f"{BASE_URL}/resetit"
        rsp = requests.post(url=resetItUrl, params={"token": token}, json={"password": password})
        if rsp.status_code == 200:
            return redirect(url_for('login'))
        else:
            print("error is here")
            return render_template('failed.html', error_code=rsp.status_code, error=rsp.json()['error'])
    return render_template("forgot_password.html", token=token)
    
@app.route('/', methods=['GET', 'POST'])
def login():
    if request.method == "POST":
        userLoginRequest = {
            "email": request.form['email'],
            "password": request.form['pass']
        }

        userLoginUrl = f"{BASE_URL}/users/login"
        rsp = requests.post(url=userLoginUrl, json=userLoginRequest)
        user_response = rsp.json()
        if rsp.status_code == 200:
            session['token'] = user_response['access_token']
            session['user_id'] = user_response['user']['id']
            session['username'] = user_response['user']['username']
            return redirect(url_for('dashboard'))
        elif rsp.status_code == 400:
            render_template('failed.html', error_code=rsp.status_code)
        elif rsp.status_code == 401:
            flash("Incorrect Password", "error")
        elif rsp.status_code == 404:
            flash("No user with this email found", "error")
        else:
            print(user_response['error'])
            flash("No user found", "error")
            # return render_template('failed.html', error_code=rsp.status_code)
    return render_template("login.html")

@app.route('/get_user/<int:id>')
def get_user(id):
    getUserUri = f"{BASE_URL}/users/{id}"
    # product_reponse = requests.get(url=f"{BASE_URL}/allproducts", headers={"Authorization": f"Bearer {session['token']}"})
    rsp = requests.get(url=getUserUri, headers={"Authorization": f"Bearer {session['token']}"})
    if rsp.status_code == 200:
        if session['user_id'] == 1:
            rspAdmin = requests.get(url=f"{BASE_URL}/users/1", headers={"Authorization": f"Bearer {session['token']}"})
            data = rsp.json()
            stock_value = data.get("stock_value")
            if stock_value is not None:
                formatted_value = "{:,}".format(stock_value)
            else:
                formatted_value = 0.00
            return render_template('user.html', user=data, admin=rspAdmin.json(), user_id=session['user_id'], ct="user", invoice_date=formatted_value)
        return render_template('user.html', user=rsp.json(), admin="none", user_id=session['user_id'], ct="user")
    elif rsp.status_code == 401:
        flash("Please login", "error")
        return redirect(url_for('login'))
    else:
        return render_template('failed.html', error_code=rsp.status_code, error=rsp.json()['error'])

@app.route('/get_product/<int:id>')
def get_product(id):
    getProductUri = f"{BASE_URL}/products/{id}"
    rsp = requests.get(url=getProductUri, headers={"Authorization": f"Bearer {session['token']}"})
    if rsp.status_code == 200:
        return render_template('get_product.html', user_data=rsp.json()) # unmarshal JSON and read data
    elif rsp.status_code == 401:
        flash("Please login", "error")
        return redirect(url_for('login'))
    else:
        return render_template('failed.html', error_code=rsp.status_code, error=rsp.json()['error'])

@app.route('/get_invoice/<int:id>')
def get_invoice(id):
    getInvoiceUri = f"{BASE_URL}/invoices/{id}"
    rsp = requests.get(url=getInvoiceUri, headers={"Authorization": f"Bearer {session['token']}"})
    print(rsp.text)
    if rsp.status_code == 200:
        return render_template('get_invoice.html', user_data=rsp.json()) # unmarshal JSON and read data
    elif rsp.status_code == 401:
        flash("Please login", "error")
        return redirect(url_for('login'))
    else:
        return render_template('failed.html', error_code=rsp.status_code, error=rsp.json()['error'])
    
@app.route('/get_receipt/<int:id>')
def get_receipt(id):
    getReceiptUri = f"{BASE_URL}/receipts/{id}"
    rsp = requests.get(url=getReceiptUri, headers={"Authorization": f"Bearer {session['token']}"})
    if rsp.status_code == 200:
        return render_template('get_receipt.html', user_data=rsp.json()) # unmarshal JSON and read data
    elif rsp.status_code == 401:
        flash("Please login", "error")
        return redirect(url_for('login'))
    else:
        return render_template('failed.html', error_code=rsp.status_code, error=rsp.json()['error'])

@app.route('/dashboard')
def dashboard():
    getUserUri = f"{BASE_URL}/users/{session['user_id']}"
    rsp = requests.get(url=getUserUri, headers={"Authorization": f"Bearer {session['token']}"})
    data = rsp.json()
    stock_value = data.get("stock_value", 0.00)
    formatted_value = "{:,}".format(stock_value)
    # if stock_value is not None:
    #     formatted_value = "{:,}".format(stock_value)
    # else:
    #     formatted_value = 0.00
    return render_template("index.html", user_id=session['user_id'], user=data, invoice_date=formatted_value)

@app.route('/list_invoices')
def list_invoices():
    listInvoicesUri = f"{BASE_URL}/invoices/admin"
    params = {'page_id': request.args.get('page_id', 1)}
    rsp = requests.get(url=listInvoicesUri, params=params, headers={"Authorization": f"Bearer {session['token']}"})
    if rsp.status_code == 500:
        flash("Please try again server error", "error")
        return render_template('failed.html', error_code=rsp.status_code, error=rsp.json()['error'])
    elif rsp.status_code == 401:
        flash("Please login", "error")
        return redirect(url_for('login'))
    return render_template("list.html", data_sent=rsp.json(), ct="invoices", user_id=session['user_id'], context="listingPagination")

@app.route('/list_receipts')
def list_receipts():
    listReceiptUrl = f"{BASE_URL}/receipts/admin"
    params = {'page_id': request.args.get('page_id', 1)}
    rsp = requests.get(url=listReceiptUrl, params=params, headers={"Authorization": f"Bearer {session['token']}"})
    if rsp.status_code == 500:
        flash("Please try again server error", "error")
        return render_template('failed.html', error_code=rsp.status_code, error=rsp.json()['error'])
    elif rsp.status_code == 401:
        flash("Please login", "error")
        return redirect(url_for('login'))
    return render_template("list.html", data_sent=rsp.json(), ct="receipts", user_id=session['user_id'], context="listingPagination")

@app.route('/list_users', methods=['GET', 'POST'])
def list_users():
    listUsersUrl = f"{BASE_URL}/users/admin"
    params = {'page_id': request.args.get('page_id', 1)}
    rsp = requests.get(url=listUsersUrl, params=params, headers={"Authorization": f"Bearer {session['token']}"})
    if rsp.status_code == 500:
        flash("Please try again server error", "error")
        return render_template('failed.html', error_code=rsp.status_code, error=rsp.json()['error'])
    elif rsp.status_code == 401:
        flash("Please login", "error")
        return redirect(url_for('login'))
    return render_template("list.html", data_sent=rsp.json(), ct="users", user_id=session['user_id'], context="listingPagination")

@app.route('/search_users', methods=['GET', 'POST'])
def search_all():
    if request.method == 'POST':
        query = request.form.get('search', request.args.get('search', 'none'))
        page_id = request.args.get('page_id', 1)
        search_context = request.args.get('search_context')
        params = {"search_query": query, "page_id": page_id, "search_context": search_context}
        listUsersUrl = f"{BASE_URL}/search/all"
        if query == "":
            if search_context == "users":
                return redirect(url_for('list_users'))
            elif search_context == "products":
                return redirect(url_for('list_products'))
            elif search_context == "receipts":
                return redirect(url_for('list_receipts'))
            else:
                return redirect(url_for('list_invoices'))
    rsp = requests.get(url=listUsersUrl, params=params, headers={"Authorization": f"Bearer {session['token']}"})
    if rsp.status_code == 500:
        flash("Please try again server error", "error")
        return render_template('failed.html', error_code=rsp.status_code, error=rsp.json()['error'])
    elif rsp.status_code == 401:
        flash("Please login", "error")
        return redirect(url_for('login'))
    if search_context == "products":
        products = requests.get(url=f"{BASE_URL}/allproducts/", headers={"Authorization": f"Bearer {session['token']}"})
        return render_template("list.html", data_sent=rsp.json(), all_products=products.json() ,ct=search_context, user_id=session['user_id'], context=f"{search_context}SearchPagination")
    return render_template("list.html", data_sent=rsp.json(), ct=search_context, user_id=session['user_id'], context=f"{search_context}SearchPagination")

# @app.route('/search_products', methods=['GET', 'POST'])
# def search_products():
#     if request.method == 'POST':
#         query = request.form.get('search')
#         listUsersUrl = f"{BASE_URL}/search/products"
#         params = {"search_word": query}
#     rsp = requests.get(url=listUsersUrl, params=params, headers={"Authorization": f"Bearer {session['token']}"})
#     if rsp.status_code == 500:
#         flash("Please try aain server error")
#         return render_template('failed.html', error_code=rsp.status_code, error=rsp.json()['error'])
#     elif rsp.status_code == 401:
#         flash("Please login")
#         return redirect(url_for('login'))
#     return render_template("list.html", data_sent=rsp.json(), ct="products", user_id=session['user_id'])

# @app.route('/search_invoices', methods=['GET', 'POST'])
# def search_invoices():
#     if request.method == 'POST':
#         query = request.form.get('search')
#         listUsersUrl = f"{BASE_URL}/search/user/invoices"
#         params = {"search_word": query}
#     rsp = requests.get(url=listUsersUrl, params=params, headers={"Authorization": f"Bearer {session['token']}"})
#     if rsp.status_code == 500:
#         flash("Please try aain server error")
#         return render_template('failed.html', error_code=rsp.status_code, error=rsp.json()['error'])
#     elif rsp.status_code == 401:
#         flash("Please login")
#         return redirect(url_for('login'))
#     return render_template("list.html", data_sent=rsp.json(), ct="products", user_id=session['user_id'])

@app.route('/list_products')
def list_products():
    listProductsUrl = f"{BASE_URL}/products"
    params = {'page_id': request.args.get('page_id', 1)}
    rsp = requests.get(url=listProductsUrl, params=params, headers={"Authorization": f"Bearer {session['token']}"})
    products = ""
    if session['user_id'] > 1:
        products = requests.get(url=f"{BASE_URL}/allproducts/", headers={"Authorization": f"Bearer {session['token']}"})
    if rsp.status_code == 500:
        flash("Please try again server error", "error")
        return render_template('failed.html', error_code=rsp.status_code, error=rsp.json()['error'])
    elif rsp.status_code == 401:
        flash("Please login", "error")
        return redirect(url_for('login'))
    if session['user_id'] > 1:
        return render_template("list.html", data_sent=rsp.json(), ct="products", all_products=products.json(), user_id=session['user_id'])
    else:
        return render_template("list.html", data_sent=rsp.json(), ct="products", all_products=products, user_id=session['user_id'])        

@app.route('/get_user_invoices/<int:id>')
def get_user_invoices(id):
    user = requests.get(url=f"{BASE_URL}/users/{id}", headers={"Authorization": f"Bearer {session['token']}"})
    params = {'page_id': request.args.get('page_id', 1)}
    # product_reponse = requests.get(url=f"{BASE_URL}/allproducts", headers={"Authorization": f"Bearer {session['token']}"})
    getUserInvoiceUrl = f"{BASE_URL}/users/invoices/{id}"
    rsp = requests.get(url=getUserInvoiceUrl, params=params, headers={"Authorization": f"Bearer {session['token']}"})
    if rsp.status_code == 200:
        data = user.json()
        stock_value = data.get("stock_value")
        if stock_value is not None:
            formatted_value = "{:,}".format(stock_value)
        else:
            formatted_value = 0.00
        if session['user_id'] == 1:
            rspAdmin = requests.get(url=f"{BASE_URL}/users/1", headers={"Authorization": f"Bearer {session['token']}"})
            adminData = rspAdmin.json()
            # stock_value = adminData.get("stock_value")
            # if stock_value is not None:
            #     formatted_value = "{:,}".format(stock_value)
            # else:
            #     formatted_value = 0.00
            return render_template("user.html", invoice=rsp.json(), user=data, admin=adminData, user_id=session['user_id'], ct='invoice', invoice_date=formatted_value)
        return render_template("user.html", invoice=rsp.json(), user=data, admin={"stock": None}, user_id=session['user_id'], ct='invoice', invoice_date=formatted_value)
    elif rsp.status_code == 401:
        flash("Please login", "error")
        return redirect(url_for('login'))
    else:
        return render_template('failed.html', error_code=rsp.status_code, error=rsp.json()['error'])
    
@app.route('/get_user_receipts/<int:id>')
def get_user_receipts(id):
    user = requests.get(url=f"{BASE_URL}/users/{id}", headers={"Authorization": f"Bearer {session['token']}"})
    params = {'page_id': request.args.get('page_id', 1)}
    getUserReceiptsUrl = f"{BASE_URL}/users/receipts/{id}"
    # product_reponse = requests.get(url=f"{BASE_URL}/allproducts", headers={"Authorization": f"Bearer {session['token']}"})
    rsp = requests.get(url=getUserReceiptsUrl, params=params, headers={"Authorization": f"Bearer {session['token']}"})
    if rsp.status_code == 200:
        data = user.json()
        stock_value = data.get("stock_value")
        if stock_value is not None:
            formatted_value = "{:,}".format(stock_value)
        else:
            formatted_value = 0.00
        if session['user_id'] == 1:
            rspAdmin = requests.get(url=f"{BASE_URL}/users/1", headers={"Authorization": f"Bearer {session['token']}"})
            adminData = rspAdmin.json()
            # stock_value = adminData.get("stock_value")
            # if stock_value is not None:
            #     formatted_value = "{:,}".format(stock_value)
            # else:
            #     formatted_value = 0.00
            return render_template("user.html", receipt=rsp.json(), user=data, admin=adminData, user_id=session['user_id'], ct='receipt', invoice_date=formatted_value)
        return render_template("user.html", receipt=rsp.json(), user=data, admin={"stock": None}, user_id=session['user_id'], ct='receipt', invoice_date=formatted_value)
    elif rsp.status_code == 401:
        flash("Please login", "error")
        return redirect(url_for('login'))
    else:
        return render_template('failed.html', error_code=rsp.status_code, error=rsp.json()['error'])

@app.route('/purchase_orders')
def purchase_orders():
    params = {'page_id': request.args.get('page_id', 1)}
    getPurchaseOrdersURL = f"{BASE_URL}/admin/purchase-orders"
    rsp = requests.get(url=getPurchaseOrdersURL, params=params, headers={"Authorization": f"Bearer {session['token']}"})
    if rsp.status_code != 200:
        return render_template('failed.html', error_code=rsp.status_code, error=rsp.json()['error'])

    rspAdmin = requests.get(url=f"{BASE_URL}/users/1", headers={"Authorization": f"Bearer {session['token']}"})
    if rspAdmin.status_code != 200:
        return render_template('failed.html', error_code=rsp.status_code, error=rsp.json()['error'])
    
    adminData = rspAdmin.json()

    return render_template("user.html", order=rsp.json(), user=adminData, admin="", user_id=session['user_id'], ct='orders', invoice_date="")

# auth.DELETE("/admin/purchase-orders/:id", server.deletePurchaseOrders)
@app.route('/purchase_orders/<string:id_param>')
def deletePurchaseOrder(id_param):
    deletePurchaseOrderUrl = f"{BASE_URL}/admin/purchase-orders/{id_param}"   
    rsp = requests.delete(url=deletePurchaseOrderUrl, headers={"Authorization": f"Bearer {session['token']}"})
    if rsp.status_code != 200:
        return render_template('failed.html', error_code=rsp.status_code, error=rsp.json()['error'])
    flash("Purchase order deleted successfully", "success")
    return redirect(url_for('purchase_orders'))

@app.route('/get_user_products/<int:id>')
def get_user_products(id):
    getUserProductsUrl = f"{BASE_URL}/users/products/{id}"
    rsp = requests.get(url=getUserProductsUrl, headers={"Authorization": f"Bearer {session['token']}"})
    if rsp.status_code == 200:
        return render_template("get_user_products.html", user_invoiceData=rsp.json())
    elif rsp.status_code == 401:
        flash("Please login", "error")
        return redirect(url_for('login'))
    else:
        return render_template('failed.html', error_code=rsp.status_code, error=rsp.json()['error'])

@app.route("/users/admin/manage/add", methods=['POST', 'GET'])
def add_admin_stock():
    if request.method == 'POST':
        id = request.form.get('productID')
        print(id)
        addAdminStockUrl = f"{BASE_URL}/users/products/admin/add/{id}"
        q = request.form.get("quantity")
        rsp = requests.post(url=addAdminStockUrl, json={"user_id": 1, "quantity": int(q)}, headers={"Authorization": f"Bearer {session['token']}"})
        if rsp.status_code == 200:
            flash("Stock Added Successfully", "success")
            return redirect(url_for('get_user', id=1))
        elif rsp.status_code == 401:
            flash("Please login", "error")
            return redirect(url_for('login'))
        else:
            return render_template('failed.html', error_code=rsp.status_code, error=rsp.json()['error'])

@app.route('/add_client_stock/<int:id>', methods=['POST', 'GET'])
def add_client_stock(id):
    if request.method == 'POST':
        quantities = request.form.getlist('quantities')
        quantities_list = [int(quantity) for quantity in quantities]
        products_id = request.form.getlist('products_id')
        products_list = [int(product_id) for product_id in products_id]
        invoice_date = request.form.get('invoiceDate')

        print(quantities_list, products_list, id)
        data = {
            "products_id": products_list,
            "quantities": quantities_list,
            "invoice_date": invoice_date
        }
        addClientStockUrl = f"{BASE_URL}/users/admin/manage/add/{id}"
        rsp = requests.post(url=addClientStockUrl, json=data, headers={"Authorization": f"Bearer {session['token']}"})
        if rsp.status_code == 200:
            flash("User Stock Added", "success")
            return redirect(url_for('get_user', id=id))
        elif rsp.status_code == 401:
            flash("Please login", "error")
            return redirect(url_for('login'))
        elif rsp.status_code == 406:
            flash(rsp.json()["error"], "error")
            return redirect(url_for('get_user', id=id))
        else:
            return render_template('failed.html', error_code=rsp.status_code)
    return render_template("add_client_stock.html")

@app.route('/users/products/sell/<int:id>', methods=['POST', 'GET'])
def reduce_client_stock(id):
    if request.method == 'POST':
        products_id = request.form.getlist('products_id')
        product_list = [int(num) for num in products_id]
        quantities = request.form.getlist('quantities')
        quantities_list = [int(num) for num in quantities]
        # amount = request.form.get('amount')

        # print(product_list, quantities_list, products_id, quantities, id)
        data = {
            "products_id": product_list,
            "quantities": quantities_list
            # "amount": int(amount)
        }

        url = f"{BASE_URL}/users/products/sell/{id}"
        rsp = requests.post(url, json=data, headers={"Authorization": f"Bearer {session['token']}"})

        if rsp.status_code == 200:
            return render_template("wait.html", user_id=id)
        elif rsp.status_code == 401:
            flash("Please login", "error")
            return redirect(url_for('login'))
        elif rsp.status_code == 406:
            flash(rsp.json()["error"], "error")
            return redirect(url_for('get_user', id=id))
        else:
            return render_template('failed.html', error_code=rsp.status_code, error=rsp.json()['error'])
    return render_template("reduce_client_stock.html")

@app.route('/users/products/sell/admin/<int:id>', methods=['POST', 'GET'])
def reduce_client_stock_by_admin(id):
    getUserUri = f"{BASE_URL}/users/{id}"
    rsp = requests.get(url=getUserUri, headers={"Authorization": f"Bearer {session['token']}"})
    if request.method == 'POST':
        products_id = request.form.getlist('products_id')
        product_list = [int(num) for num in products_id]
        quantities = request.form.getlist('quantities')
        quantities_list = [int(num) for num in quantities]
        phone_number = request.form.get('phone_number')
        mpese_receipt_number = request.form.get('mpesa_receipt_number')
        description = request.form.get('description')
        amount = request.form.get('amount')

        # print(product_list, quantities_list, products_id, quantities, id)
        data = {
            "products_id": product_list,
            "quantities": quantities_list,
            "phone_number": phone_number,
            "mpesa_receipt_number": mpese_receipt_number,
            "description": description,
            "amount": int(amount)
        }

        url = f"{BASE_URL}/users/admin/reduce_client_stock/{id}"
        rsp = requests.post(url, json=data, headers={"Authorization": f"Bearer {session['token']}"})

        if rsp.status_code == 200:
            flash("User 3rd Party Payment Recorded", "success")
            return redirect(url_for('get_user', id=id))
        elif rsp.status_code == 401:
            flash("Please login", "error")
            return redirect(url_for('login'))
        elif rsp.status_code == 406:
            flash(rsp.json()["error"], "error")
            return redirect(url_for('reduce_client_stock_by_admin', id=id))
        else:
            return render_template('failed.html', error_code=rsp.status_code, error=rsp.json()['error'])
    return render_template("admin_reduce_stock.html", user=rsp.json())

@app.route("/search/transactions", methods=['POST', 'GET'])
def search_transactions():
    if request.method == "POST":
        query = request.form.get("search", request.args.get("search", "none"))
        page_id = request.args.get("page_id", 1)
        search_context = request.args.get("search_context")
        params = {"search_query": query, "page_id": page_id, "search_context": search_context}
        listUsersUrl = f"{BASE_URL}/search/all"
    rsp = requests.get(url=listUsersUrl, params=params, headers={"Authorization": f"Bearer {session['token']}"})
    if rsp.status_code == 500:
        flash("Please try again server error", "error")
        return render_template('failed.html', error_code=rsp.status_code, error=rsp.json()['error'])
    elif rsp.status_code == 401:
        flash("Please login", "error")
        return redirect(url_for('login'))
    return render_template("transactions.html", data_sent=rsp.json(), user_id=session['user_id'], action=f"search_{search_context}")

@app.route("/transactions")
def getAllTransactions():
    params = {'page_id': request.args.get('page_id', 1)}
    url = f"{BASE_URL}/transactions/all"
    rsp = requests.get(url, params=params, headers={"Authorization": f"Bearer {session['token']}"})
    return render_template("transactions.html", id=0, data_sent=rsp.json(), user_id=session['user_id'], action="all_transactions")

@app.route("/transactions/successfull")
def getSuccessfulTransactions():
    params = {'page_id': request.args.get('page_id', 1)}
    url = f"{BASE_URL}/transactions/successfull"
    rsp = requests.get(url, params=params, headers={"Authorization": f"Bearer {session['token']}"})
    return render_template("transactions.html", id=0, data_sent=rsp.json(), user_id=session['user_id'], action="successful_transactions")

@app.route("/transactions/failed")
def getFailedTransactions():
    params = {'page_id': request.args.get('page_id', 1)}
    url = f"{BASE_URL}/transactions/failed"
    rsp = requests.get(url, params=params, headers={"Authorization": f"Bearer {session['token']}"})
    return render_template("transactions.html", id=0, data_sent=rsp.json(), user_id=session['user_id'], action="failed_transactions")

@app.route("/transactions/users/<int:user_id>")
def getUserAllTransactions(user_id):
    params = {'page_id': request.args.get('page_id', 1)}
    url = f"{BASE_URL}/user/transactions/all/{user_id}"
    rsp = requests.get(url, params=params, headers={"Authorization": f"Bearer {session['token']}"})
    return render_template("transactions.html", id=session["user_id"], data_sent=rsp.json(), user_id=user_id, action="user_all_transactions")

@app.route("/transactions/users/successful/<int:user_id>")
def getUserSuccessfulTransactions(user_id):
    params = {'page_id': request.args.get('page_id', 1)}
    url = f"{BASE_URL}/user/transactions/successful/{user_id}"
    rsp = requests.get(url, params=params, headers={"Authorization": f"Bearer {session['token']}"})
    return render_template("transactions.html", id=session["user_id"], data_sent=rsp.json(), user_id=user_id, action="user_successful_transactions")

@app.route("/transactions/users/failed/<int:user_id>")
def getUserFailedTransactions(user_id):
    params = {'page_id': request.args.get('page_id', 1)}
    url = f"{BASE_URL}/user/transactions/failed/{user_id}"
    rsp = requests.get(url, params=params, headers={"Authorization": f"Bearer {session['token']}"})
    return render_template("transactions.html", id=session["user_id"], data_sent=rsp.json(), user_id=user_id, action="user_failed_transactions")

# @app.route("/transactions/<transaction_number>", methods=['POST', 'GET'])
# def getTransaction(transaction_number):
#     params = {'page_id': request.args.get('page_id', 1)}
#     url = f"{BASE_URL}/transactions/all"
#     rsp = requests.get(url, params=params, headers={"Authorization": f"Bearer {session['token']}"})
#     return render_template("transactions.html", data_sent=rsp.json(), user_id=session['user_id'], action="all_transactions")

@app.route("/download/invoice/<string:id_param>", methods=['POST', 'GET'])
def invoiceDownload(id_param):
    url = f"{BASE_URL}/invoice/download/{id_param}"
    response = requests.get(url=url, headers={"Authorization": f"Bearer {session['token']}"})
    data = response.json()

    pdf_bytes = base64.b64decode(data['invoice_pdf'])

    if response.status_code == 200:
        # flash("Invoice Downloaded Successfully", "success")
        return send_file(BytesIO(pdf_bytes), as_attachment=True, mimetype='application/pdf', download_name=f"INV-{id_param}.pdf")
    elif response.status_code == 401:
        flash("Please login", "error")
        return redirect(url_for('login'))
    else:
        return render_template('failed.html', error_code=response.status_code, error=response.json()['error'])
    
@app.route("/download/receipt/<string:id_param>", methods=['POST', 'GET'])
def receiptDownload(id_param):
    url = f"{BASE_URL}/receipt/download/{id_param}"
    response = requests.get(url=url, headers={"Authorization": f"Bearer {session['token']}"})
    data = response.json()

    if response.status_code == 200:
        pdf_bytes = base64.b64decode(data['receipt_pdf'])
        # flash("Receipt Downloaded Successfully", "success")
        return send_file(BytesIO(pdf_bytes), as_attachment=True, mimetype='application/pdf', download_name=f"RCPT-{id_param}.pdf")
    elif response.status_code == 401:
        flash("Please login", "error")
        return redirect(url_for('login'))
    else:
        return render_template('failed.html', error_code=response.status_code, error=response.json()['error'])
    
    # auth.GET("/admin/purchase-orders/:id", server.downloadPurchaseOrders)
@app.route("/download/purchase-order/<string:id_param>", methods=['POST', 'GET'])
def purchaseOrderDownload(id_param):
    url = f"{BASE_URL}/admin/purchase-orders/{id_param}"
    response = requests.get(url=url, headers={"Authorization": f"Bearer {session['token']}"})
    data = response.json()

    if response.status_code == 200:
        pdf_bytes = base64.b64decode(data['purchase_order_pdf'])
        return send_file(BytesIO(pdf_bytes), as_attachment=True, mimetype='application/pdf', download_name=f"PO-{id_param}.pdf")
    elif response.status_code == 401:
        flash("Please login", "error")
        return redirect(url_for('login'))
    else:
        return render_template('failed.html', error_code=response.status_code, error=response.json()['error'])

@app.route("/download/statement/<int:id>", methods=['POST', 'GET'])
def statement_download(id):
    url = f"{BASE_URL}/statements/{id}"
    response = requests.get(url=url, headers={"Authorization": f"Bearer {session['token']}"})
    data = response.json()

    if response.status_code == 200:
        pdf_bytes = base64.b64decode(data['statement_pdf'])
        # flash("Receipt Downloaded Successfully", "success")
        return send_file(BytesIO(pdf_bytes), as_attachment=True, mimetype='application/pdf', download_name=f"statement.pdf")
    elif response.status_code == 202:
        flash("The user has no transactions history yet", "info")
        return redirect(url_for("get_user", id=id))
    else:
        flash("Failed to download statement", "error")
        return redirect(url_for("get_user", id=id))
    
@app.route("/download/reports", methods=['POST', 'GET'])
def report_download():
    selected_option = request.form.get('products_id')
    reportType = "admin"
    if selected_option == "0":
        reportType = "users"
    from_date =request.form.get("from_date")
    to_date = request.form.get("to_date")
    params = {"start_date": from_date, "end_date": to_date}
    url = f"{BASE_URL}/admin/{reportType}_reports"
    response = requests.post(url=url, json=params, headers={"Authorization": f"Bearer {session['token']}"})
    data = response.json()
    if response.status_code == 200:
        excel_bytes = base64.b64decode(data['data'])
        return send_file(BytesIO(excel_bytes), as_attachment=True, mimetype='application/vnd.openxmlformats-officedocument.spreadsheetml.sheet', download_name=f"report.xlsx")
    else:
        flash("Failed to download statement", "error")
        return redirect(url_for("get_user", id=id))


@app.route("/purchase_order")
def create_purchase_order():
    rspAdmin = requests.get(url=f"{BASE_URL}/users/1", headers={"Authorization": f"Bearer {session['token']}"})
    adminData = rspAdmin.json()
    return render_template("purchase-order.html", admin=adminData)

@app.route("/download/purchase-order", methods=['POST'])
def purchase_order_download():
    form_data = request.form.get('data')
    if not form_data:
        flash("No data provided", "error")
        return redirect(url_for("get_user", id=1))
    
    try:
        body = json.loads(form_data)
    except ValueError:
        flash("Invalid data format", "error")
        return redirect(url_for("get_user", id=1))

    url = f"{BASE_URL}/admin/purchase-order"
    response = requests.post(url=url, json=body, headers={"Authorization": f"Bearer {session['token']}"})
    data = response.json()

    if response.status_code == 200:
        pdf_bytes = base64.b64decode(data['purchase_order_pdf'])
        return send_file(BytesIO(pdf_bytes), as_attachment=True, mimetype='application/pdf', download_name="purchase_order.pdf")
    elif response.status_code == 202:
        return redirect(url_for("get_user", id=1))
    else:
        flash("Failed to download statement", "error")
        return redirect(url_for("get_user", id=1))

    # url = f"{BASE_URL}/admin/purchase-order"
    # body = request.json
    # response = requests.post(url=url, json=body, headers={"Authorization": f"Bearer {session['token']}"})
    # data = response.json()

    # if response.status_code == 200:
    #     pdf_bytes = base64.b64decode(data['statement_pdf'])
    #     return send_file(BytesIO(pdf_bytes), as_attachment=True, mimetype='application/pdf', download_name="purchase_order.pdf")
    # elif response.status_code == 202:
    #     return redirect(url_for("get_user", id=1))
    # else:
    #     flash("Failed to download statement", "error")
    #     return redirect(url_for("get_user", id=1))

@app.route("/request_stock/<int:id>", methods=["POST", "GET"])
def request_stock(id):
        if request.method == 'POST':
            products_id = request.form.getlist('products_id_request')
            product_list = [int(num) for num in products_id]
            quantities = request.form.getlist('quantities')
            quantities_list = [int(num) for num in quantities]

            data = {
                "products": product_list,
                "quantities": quantities_list
            }
            url = f"{BASE_URL}/users/request_stock/{id}"
            rsp = requests.post(url, json=data, headers={"Authorization": f"Bearer {session['token']}"})

            if rsp.status_code == 200:
                flash("Request Recorded", "success")
                return redirect(url_for('list_products'))  # Redirect to some success page
            elif rsp.status_code == 401:
                flash("Please login", "error")
                return redirect(url_for('login'))
            else:
                return render_template('failed.html', error_code=rsp.status_code, error=rsp.json()['error'])
        return render_template("failed.html")

@app.route("/history/received/<int:user_id>")
def user_received_history(user_id):
        if user_id == 1:
            url = f"{BASE_URL}/history/admin"
        else:
            url = f"{BASE_URL}/history/received/{user_id}"
        rsp = requests.get(url, headers={"Authorization": f"Bearer {session['token']}"})
        if rsp.status_code == 200:
            myData = rsp.json()
            if myData:
                if session['user_id'] == 1:
                    sorted_data = sorted(myData, key=lambda x: x['issued_date'], reverse=True)
                    return render_template('admin-history.html', user_id=session['user_id'], data_sent=sorted_data, action="received")
                
                sorted_data = sorted(myData.items(), key=lambda x: x[0], reverse=True)
                sorted_dict = OrderedDict(sorted_data)
                return render_template('history.html', user_id=session['user_id'], data_sent=sorted_dict, action="received")
            else:
                flash("No data found", "info")
                return redirect(url_for('dashboard'))
        elif rsp.status_code == 401:
            flash("Please login", "error")
            return redirect(url_for('login'))
        else:
            return render_template('failed.html', error_code=rsp.status_code, error=rsp.json()['error'])
        
@app.route("/history/sold/<int:user_id>")
def user_sold_history(user_id):
        if user_id == 1:
            url = f"{BASE_URL}/history/all_received"
        else:
            url = f"{BASE_URL}/history/sold/{user_id}"
        rsp = requests.get(url, headers={"Authorization": f"Bearer {session['token']}"})
        if rsp.status_code == 200:
            if session['user_id'] == 1:
                sorted_data = sorted(rsp.json().items(), key=lambda x: x[0], reverse=True)
                sorted_dict = OrderedDict(sorted_data)
                return render_template('admin-history.html', user_id=session['user_id'], data_sent=sorted_dict, action="sold")
            return render_template('history.html', user_id=session['user_id'], data_sent=rsp.json(), action="sold")
        elif rsp.status_code == 401:
            flash("Please login", "error")
            return redirect(url_for('login'))
        else:
            return render_template('failed.html', error_code=rsp.status_code, error=rsp.json()['error'])
        
@app.route("/history/debt/<int:user_id>")
def user_debt_history(user_id):
        if user_id == 1:
            url = f"{BASE_URL}/history/all_debt"
        else:    
            url = f"{BASE_URL}/history/debt/{user_id}"
        rsp = requests.get(url, headers={"Authorization": f"Bearer {session['token']}"})
        if rsp.status_code == 200:
            if session['user_id'] == 1:
                data = rsp.json()
                user_total_quantity = {}
                user_total_price = {}
                if data is not None:
                    for entry in data:
                        if entry['Data'] is not None:
                            user = entry['user']
                            total_price = sum(product['price'] for product in entry['Data'])
                            total_quantity = sum(product['quantity'] for product in entry['Data'])
                            user_total_price[user] = total_price
                            user_total_quantity[user] = total_quantity
                return render_template('admin-history.html', user_id=session['user_id'], price=user_total_price, quantity=user_total_quantity, data_sent=rsp.json(), action="debt")
            return render_template('history.html', user_id=session['user_id'], data_sent=rsp.json(), action="debt")
        elif rsp.status_code == 401:
            flash("Please login", "error")
            return redirect(url_for('login'))
        else:
            return render_template('failed.html', error_code=rsp.status_code, error=rsp.json()['error'])

# @app.route('/users/products/sell/<int:id>', methods=['POST', 'GET'])
# def reduce_client_stock(id):
#     if request.method == 'POST':
#         products_id = request.form.getlist('products_id')
#         product_list = [int(num) for num in products_id]
#         quantities = request.form.getlist('quantities')
#         quantities_list = [int(num) for num in quantities]

#         print(product_list, quantities_list, products_id, quantities, id)
#         data = {
#             "products_id": product_list,
#             "quantities": quantities_list
#         }

#         url = f"{BASE_URL}/users/products/sell/{id}"
#         rsp = requests.post(url, json=data, headers={"Authorization": f"Bearer {session['token']}"})

#         if rsp.status_code == 200:
#             # flash("STK push sent")
#             return redirect(url_for('get_user', id=id))  # Redirect to some success page
#         elif rsp.status_code == 401:
#             flash("Please login")
#             return redirect(url_for('login'))
#         else:
#             return render_template('failed.html', error_code=rsp.status_code, error=rsp.json()['error'])
#     return render_template("reduce_client_stock.html")

@app.route("/notify", methods=["POST", "GET"])
def notify():
    if request.method == "POST":
        data = request.get_json()
        return render_template("failed.html", error_code=data.get('transactionID'),error=data.get('status'))

@app.errorhandler(ConnectionError)
def handle_connection_error(error):
    return render_template('failed.html', error_code=500, error=str(error), connection=True)

@app.errorhandler(InternalServerError)
def handle_server_error(error):
    return render_template('failed.html', error_code=401, error=str(error), connection=False)

if __name__ == '__main__':
    app.run(debug=True)
