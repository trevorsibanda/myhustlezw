import React, { Component } from "react";
import { Link, NavLink, Switch, Route, Redirect, useHistory } from "react-router-dom"
import v1 from "../api/v1";
import AuthTopNav from "../components/AuthTopNav";
import LoginPage from "../components/LoginPage";

const countryCodes = require('../api/countryCodes.json')



class SignupPage extends Component {

    constructor(props){
        super(props)
        v1.page.set({title: 'Create your account'})

        this.state = {
            form: {
                fullname: '',
                username: '',
                email: '',
                country_code: '+263',
                phone_number: '',
                role: 'content creator',
            },
            loggedIn: false,
            usernameAvailable: false,
            usernameCheckLoading: false,
        }



        this.signup = () => {
            let form = this.state.form
            if(form.fullname.length < 3) {
                return alert('Fullname should be at least 3 characters long')
            }

            if(form.username.length < 1) {
                return alert('Username cannot be empty.')
            }

            if(form.phone_number.length < 4) {
                return alert('Please enter a valid phone number')
            }

            if(form.email.length < 2) {
                return alert('Please enter a valid email')
            }

            if(form.role === 'select') {
                form.role = 'content creator'
            }

            if (form.password.length < 6) {
                return alert('Password should be at least 6 characters long')
            }

            if (form.password !== form.password_confirm) {
                return alert('Passwords do not match.')
            }

            //if(! this.state.usernameAvailable) {
            //    return alert('<p>Username <b>'+form.username+'</b> is not available</p>')
            //}

            v1.security.signup(form, 'creator').then(response => {
                if (response.error) {
                    alert({
                        title: 'Failed to create your account',
                        text: response.error,
                        icon: 'warning',
                        timer: 10000,
                    })
                    this.setState({ loading: false })
                } else if (response._id) {
                    alert({
                        text: 'Welcome to MyHustle! ' + response.fullname,
                        timer: 1500,
                        toast: true
                    })
                    if (this.props.reload)
                        this.props.reload(response)
                    this.setState({ loggedIn: true })
                }
            }).catch(_ => {
                this.setState({ loading: false })
            })
        }

        v1.page.track()
    } 

    render(){
        return this.state.loggedIn ? <Redirect to="/creator/getting-started" /> : (
            <>
                <div class="box">
                    <div class="box-header">
                        <h4 class="box-title">Create a new account.</h4>
                    </div>
                </div>
                    <form action="#">
                        <div class="row">


                            <div class="input-group col-lg-6 mb-4">
                                <div class="input-group-prepend">
                                    <span class="input-group-text bg-white px-4 border-md border-right-0">
                                        <i class="fa fa-user text-muted"></i>
                                    </span>
                                </div>
                                <input id="fullname" type="text" value={this.state.fullname} onChange={evt =>  { let form = this.state.form; form.fullname = evt.target.value; this.setState({form}); }}name="fullname" placeholder="Fullname" class="form-control bg-white border-left-0 border-md" />
                            </div>


                            <div class="input-group col-lg-6 mb-4">
                                <div class="input-group-prepend">
                                    <span class="input-group-text bg-white px-4 border-md border-right-0">
                                        <i class={this.state.usernameCheckLoading ? "fa fa-loading text-muted ": this.state.usernameAvailable ? "fa fa-check text-muted" : "fa fa-cancel text-muted" }></i>
                                    </span>
                                </div>
                                <input id="username" type="text" name="username" value={this.state.username} onChange={evt =>  { let form = this.state.form; form.username = evt.target.value; this.setState({form}); }}placeholder="Username" class="form-control bg-white border-left-0 border-md" />
                            </div>

                            <div class="input-group col-lg-12 mb-4">
                                <div class="input-group-prepend">
                                    <span class="input-group-text bg-white px-4 border-md border-right-0">
                                        <i class="fa fa-phone-square text-muted"></i>
                                    </span>
                                </div>
                                <select id="countryCode" value={this.state.country_code} onChange={evt =>  { let form = this.state.form; form.country_code = evt.target.value; this.setState({form}); }}name="countryCode" style={{ 'maxWidth': "80px" }} class="custom-select form-control bg-white border-left-0 border-md font-weight-bold text-muted">
                                    { countryCodes.map( (cc) =>{
                                        return <option value={cc.dial_code}>{cc.dial_code} {cc.name}</option>
                                    })
                                    }
                                </select>
                                <input id="phoneNumber" type="tel" name="phone" value={this.state.phone_number} onChange={evt =>  { let form = this.state.form; form.phone_number = evt.target.value; this.setState({form}); }}placeholder="Phone Number" class="form-control bg-white border-md border-left-0 pl-3" />
                            </div>.


                                    <div class="input-group col-lg-12 mb-4">
                                <div class="input-group-prepend">
                                    <span class="input-group-text bg-white px-4 border-md border-right-0">
                                        <i class="fa fa-envelope text-muted"></i>
                                    </span>
                                </div>
                                <input id="email" type="email" name="email" value={this.state.email} onChange={evt =>  { let form = this.state.form; form.email = evt.target.value; this.setState({form}); }}placeholder="Your Email Address" class="form-control bg-white border-left-0 border-md" />
                            </div>

                            <div class="input-group col-lg-12 mb-4">
                                <div class="input-group-prepend">
                                    <span class="input-group-text bg-white px-4 border-md border-right-0">
                                        <i class="fa fa-black-tie text-muted"></i>
                                    </span>
                                </div>
                                <select id="job" name="jobtitle" value={this.state.role} onChange={evt =>  { let form = this.state.form; form.role = evt.target.value; this.setState({form}); }}class="form-control custom-select bg-white border-left-0 border-md">
                                    <option value="select">I best identify as a ...</option>
                                    <option value="content creator">Content creator </option>
                                    <option value="professional model">Professional model</option>
                                    <option value="instagram influencer">Instagram Influencer</option>
                                    <option value="entertainer">Entertainer</option>
                                    <option value="youtuber">Youtuber</option>
                                    <option value="fitness trainer">Fitness Trainer</option>
                                    <option value="teacher/educator">Teacher/Educator</option>
                                    <option value="hustler">Hustler</option>
                                    <option value="comedian">Technician</option>
                                </select>
                            </div>
                            <div class="input-group col-lg-6 mb-4">
                                <div class="input-group-prepend">
                                    <span class="input-group-text bg-white px-4 border-md border-right-0">
                                        <i class="fa fa-lock text-muted"></i>
                                    </span>
                                </div>
                                <input id="password" type="password" value={this.state.form.password} onChange={evt =>  { let form = this.state.form; form.password = evt.target.value; this.setState({form}); }}name="password" placeholder="Password" class="form-control bg-white border-left-0 border-md" />
                            </div>


                            <div class="input-group col-lg-6 mb-4">
                                <div class="input-group-prepend">
                                    <span class="input-group-text bg-white px-4 border-md border-right-0">
                                        <i class="fa fa-lock text-muted"></i>
                                    </span>
                                </div>
                                <input id="passwordConfirmation" value={this.state.form.password_confirm} onChange={evt =>  { let form = this.state.form; form.password_confirm = evt.target.value; this.setState({form}); }} type="password" name="passwordConfirmation" placeholder="Confirm Password" class="form-control bg-white border-left-0 border-md" />
                            </div>


                            <div class="form-group col-lg-12 mx-auto mb-0">
                                <btn onClick={this.signup} class="btn btn-info btn-block py-2">
                                    <span class="font-weight-bold">Create your account</span>
                                </btn>
                            </div>



                            <div class="text-center padding-top-10">
                                <p class="text-muted font-weight-bold">Already Registered? <Link to="/auth/login" class="text-primary ml-2">Login</Link></p>
                            </div>

                        </div>
                    </form>
            </>
        )
    }
}


class RecoverPassword extends Component {
    constructor(props) {
        super(props)
        v1.page.set({title: 'Recover your password'})

        this.state = {
            sentOTP: false, verifiedOTP: false, passwordChanged: false, newPassword: "", email: '', otp: ''
        }
        this.resendCode = () => {
            v1.security.resetPasswordRequest({email: this.state.email}).then(resp => {
                if (resp.status === 'ok') {
                    alert('Resent the verification code.')
                } else {
                    alert(resp.error)
                }
            }).catch(alert)
        }

        this.processResetRequest = () => {
            v1.security.processPasswordReset({
                email: this.state.email,
                otp: this.state.otp,
                newPassword: this.state.newPassword,
            }).then(resp => {
                if (resp.status === 'changed') {
                    alert('Password changed successfully. You can now login with your new password.')
                    useHistory().push('/auth/login')
                } else {
                    alert(resp.error)
                }
            }).catch(err => {
                alert(err.error)
            })
        }

        this.requestReset = () => {
            if (!v1.util.validateEmail(this.state.email)) {
                alert('Please enter a valid email address.')
                return
            }
            v1.security.resetPasswordRequest({ email: this.state.email }).then(resp => {
                if (resp.status === 'ok') {
                    alert("Sent the verification code. Check your email.");
                    this.setState({sentOTP: true, showemail: false })
                } else {
                    alert(resp.error)
                }
            }).catch(err => {
                alert(err.error)
            })
        }

        v1.page.track()
    }

    render() {
        return (
            <>
                <div class="box">
                    <div class="box-header">
                        <h4 class="box-title">Let's help you reset your password.</h4>
                    </div>
                </div>
                <div class="row justify-content-center padding-top-50">
                    <div class="col-12  ">
                        <div class="card py-5 px-3">
                            
                            {
                                this.state.sentOTP ?
                                    <>
                                        <h5 class="m-0 padding-bottom-40">Enter the verification code you received by email.</h5>
                                        <span class="mobile-text">We sent a verification code to your email address.</span>

                                        <p>We sent a password reset code to <b class="text-danger">{this.state.email}</b>. Enter it below and click verify.</p>
                                        <div class="d-flex flex-row mt-5" data-children-count="4">
                                            <input type="text"  class="form-control" autofocus="" value={this.state.otp} onChange={(evt) => { this.setState({ otp: evt.target.value }) }} placeholder="6 digit OTP sent to your email" />
                                        </div>
                                        <label>* New password</label>
                                        <input type="text" class="form-control" placeholder="Your new password" value={this.state.newPassword} onChange={(evt) => { this.setState({ newPassword: evt.target.value }) }} />
                                        <p><small>Password is shown in clear, min of 6 characters. Choose a strong password.</small></p>
                                        
                                        <button class="btn btn-block btn-primary" onClick={this.processResetRequest}><i class="fa fa-check"></i> Verify code</button>
                                    </> :
                                    <>
                                    <h5 class="m-0 padding-bottom-40">Enter the email you used to signup.</h5>
                                    <div class="form-group">
                                        <label>* Email address</label>
                                        <input type="email" class="form-control" placeholder="myemail@gmail.com" value={this.state.email} onChange={(evt) => { this.setState({ email: evt.target.value }) }} />
                                        <p><small>Email you used to register your account.</small></p>
                                        <button class="btn btn-block btn-info" onClick={this.requestReset}><i class="fa fa-arrow-right"></i> Send password reset code</button>
                                    </div>
                                        
                                    </>
                            }
                            
                        </div>
                    </div>
                </div>
            </>
        )
    }
} 


class AuthPage extends Component {
    constructor(props) {
        super(props)
        this.state = {
            loggedIn: false,
        }

        v1.user.loggedIn().then(user => {
            console.log(user)
            if(user)
            this.setState({ user, loggedIn: true })
        })

    }

    render() {
        return this.state.loggedIn ? <Redirect to="/creator/dashboard" /> : (
            <>
                <AuthTopNav user={this.props.user} />
                <div class="container">
                    <div class="row padding-top-40 mt-4 align-items-center">
                        <div class="col-md-5 pr-lg-5 mb-5 mb-md-0">
                            
                        </div>
                        <div class="col-md-7 col-lg-6 ml-auto">
                            <ul class="nav nav-pills justify-content-end margin-bottom">
                                <li class=" nav-item"> <NavLink className="nav-link" to="/auth/login">Login</NavLink> </li>
                                <li class="nav-item"> <NavLink className="nav-link" to="/auth/signup">Signup</NavLink> </li>
                                <li class="nav-item"> <NavLink className="nav-link" to="/auth/recover_password">Forgot password?</NavLink> </li>
                            </ul>
                        </div>
                    </div>
                    <div class="row  padding-bottom-150">
                        <div class="col-md-5 pr-lg-5 mb-5 mb-md-0">
                            <img src="https://res.cloudinary.com/mhmd/image/upload/v1569543678/form_d9sh6m.svg" alt="" class="img-fluid mb-3 d-none d-md-block" />
                            <h1>MyHustle</h1>
                            
                        </div>
                        <div class="col-md-7 col-lg-6 ml-auto">
                            <Switch>
                                <Route path="/auth/login" render={props => <LoginPage reload={this.props.reload} />} />
                                <Route path="/auth/signup" render={props => <SignupPage reload={this.props.reload} />} />
                                <Route path="/auth/recover_password" component={RecoverPassword} />
                                <Route path="/">
                                    <Redirect to="/auth/login" />
                                </Route>
                            </Switch>
                        </div>
                    </div>
                </div>

            </>
        );
    }
}

export default AuthPage;

