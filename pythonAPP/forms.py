from flask_wtf import FlaskForm
from wtforms import StringField, PasswordField, SubmitField, IntegerField
from wtforms.validators import DataRequired, Length

class CreateUserForm(FlaskForm):
    username = StringField('Username', validators=[DataRequired()])
    password = PasswordField('Password', validators=[DataRequired(), Length(min=6)])
    email = StringField('Email', validators=[DataRequired()])
    phoneNumber = StringField('PhoneNumber', validators=[DataRequired()])
    address = StringField('Address', validators=[DataRequired()])
    submit = SubmitField('Login')

class EditProductForm(FlaskForm):
    productName = StringField('ProductName', validators=[DataRequired()])
    unitPrice = IntegerField('UnitPrice', validators=[DataRequired()])
    packSize = StringField('Packsize', validators=[DataRequired()])
    submit = SubmitField('Login')

class CreateProductForm(FlaskForm):
    productName = StringField('ProductName', validators=[DataRequired()])
    unitPrice = IntegerField('UnitPrice', validators=[DataRequired()])
    packsize = StringField('Packsize', validators=[DataRequired()])
    submit = SubmitField('Login')

class ChangePasswordForm(FlaskForm):
    oldPassword = StringField('OldPassword', validators=[DataRequired()])
    newPassword = PasswordField('NewPassword', validators=[DataRequired(), Length(min=6)])
    submit = SubmitField('Login')

class ManageUserForm(FlaskForm):
    email = StringField('Email', validators=[DataRequired()])
    phoneNumber = StringField('PhoneNumber', validators=[DataRequired()])
    address = StringField('Address', validators=[DataRequired()])
    username = StringField('Username', validators=[DataRequired()])
    submit = SubmitField('Login')

class ResetPasswordForm(FlaskForm):
    email = StringField('Email', validators=[DataRequired()])
    submit = SubmitField('Login')

class ResetItForm(FlaskForm):
    Password = PasswordField('Password', validators=[DataRequired(), Length(min=6)])
    submit = SubmitField('Login')

class LoginForm(FlaskForm):
    email = StringField('Email', validators=[DataRequired()])
    password = PasswordField('Password', validators=[DataRequired(), Length(min=6)])
    submit = SubmitField('Login')

class AddAdminStockQuantity(FlaskForm):
    quantity = IntegerField('quantity', validators=[DataRequired()])
    submit = SubmitField('Login')
