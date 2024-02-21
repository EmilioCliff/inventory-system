from flask import Flask, render_template
from forms import CreateUserForm, CreateProductForm

app = Flask(__name__)

app.config['SECRET_KEY'] = "32e234353t4rffbfbfgxx"

data = [
    {
    "username": "Niko Hapa",
    "email": "email@gmail.com",
    "phone_number": "0912323423",
    "address": "myadress",
    "stock": [ { 
        "product_name": "Generally",
        "unit_price": 400,
        "packsize": "50 Testkit"
        },
            {
        "product_name": "Generally",
        "unit_price": 400,
        "packsize": "50 Testkit"
        }
    ]
            },
    {
    "username": "Niko Hapa",
    "email": "email@gmail.com",
    "phone_number": "0912323423",
    "address": "myadress",
    "stock": [ {
        "product_name": "Generally",
        "unit_price": 400,
        "packsize": "50 Testkit"
    },
        {
        "product_name": "Generally",
        "unit_price": 400,
        "packsize": "50 Testkit"
        }
            ]
                },
    {
    "username": "Niko Hapa",
    "email": "email@gmail.com",
    "phone_number": "0912323423",
    "address": "myadress"
    },
     {
    "username": "Niko Hapa",
    "email": "email@gmail.com",
    "phone_number": "0912323423",
    "address": "myadress"
    }
]

pdata = [
    {
        "product_name": "Generally",
        "unit_price": 400,
        "packsize": "50 Testkit"
    },
        {
        "product_name": "Generally",
        "unit_price": 400,
        "packsize": "50 Testkit"
    },
        {
        "product_name": "Generally",
        "unit_price": 400,
        "packsize": "50 Testkit"
    },
        {
        "product_name": "Generally",
        "unit_price": 400,
        "packsize": "50 Testkit"
    },
        {
        "product_name": "Generally",
        "unit_price": 400,
        "packsize": "50 Testkit"
    },
        {
        "product_name": "Generally",
        "unit_price": 400,
        "packsize": "50 Testkit"
    }
]

rdata = [
    {
        "receipt_number": "989",
        "user_receipt_username": "Emilio CLiff",
        "created_at": "12 Feb 2024"
    },
       {
        "receipt_number": "989",
        "user_receipt_username": "Emilio CLiff",
        "created_at": "12 Feb 2024"
    },
       {
        "receipt_number": "989",
        "user_receipt_username": "Emilio CLiff",
        "created_at": "12 Feb 2024"
    },
       {
        "receipt_number": "989",
        "user_receipt_username": "Emilio CLiff",
        "created_at": "12 Feb 2024"
    },
       {
        "receipt_number": "989",
        "user_receipt_username": "Emilio CLiff",
        "created_at": "12 Feb 2024"
    },
       {
        "receipt_number": "989",
        "user_receipt_username": "Emilio CLiff",
        "created_at": "12 Feb 2024"
    }
]

user = {
    "username": "Niko Hapa",
    "email": "email@gmail.com",
    "phone_number": "0912323423",
    "address": "myadress",
    "stock": [ {
        "product_name": "Generally 1",
        "product_quantity": 10
    },
        {
        "product_name": "Generally 2",
        "product_quantity": 100
    }]
}

@app.route('/')
def index():
    return render_template('list.html', mdata="products", data_sent=pdata, real=user)

@app.route('/create', methods=['POST', 'GET'])
def create():
    form= CreateUserForm()
    if form.validate_on_submit():
        print("%s, %s, %s, %s, %s", form.username.data, form.email.data, form.password.data, form.address.data, form.phoneNumber.data)
    return render_template('create.html', form=form)

@app.route('/createee', methods=['POST', 'GET'])
def createee():
    form= CreateProductForm()
    if form.validate_on_submit():
        print("success")
    return render_template('createee.html', form=form)

if __name__ == '__main__':
    app.run(debug=True)