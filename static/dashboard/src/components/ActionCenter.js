import React, { Component } from "react";
import { Link } from "react-router-dom"
import v1 from "../api/v1";


class VerifyPhoneNumberAction extends Component{
    constructor(props){
        super(props)

        this.state = {verified: this.props.user.phoneVerified, showPhone: false, enableVerify: true, phone: this.props.user.phone_number, otp: ''}
        this.resendCode = () => {
            v1.security.resendSMSCode().then(resp => {
                if(resp.status === 'sent'){
                    alert('Resent the verification code.')
                }else{
                    alert(resp.error)
                }
            }).catch(alert)
        }

        this.verifyPhone = () => {
            v1.security.verifyPhone(this.state.otp).then(resp => {
                if (resp.status === 'verified') {

                    window.location.reload()
                }else {
                    alert(resp.error)
                }
            }).catch(err => {
                alert(err.error)
            })
        }

        this.updatePhone = () => {
            v1.security.updatePhone(this.state.phone).then(resp => {
                if (resp.status === 'updated'){
                    this.setState({enableVerify: true, showPhone: false})
                }else {
                    alert(resp.error)
                }
            }).catch(err => {
                alert(err.error)
            })
        }
    }

    render(){
        return (
            <>
                <div class="box">
                    <div class="box-header">
                        <h4 class="box-title">Please verify your phone number to activate your account.</h4>
                    </div>
                </div>
                <div class="row justify-content-center padding-top-50">
                    <div class="col-md-6 col-xs-12 d-flex justify-content-center align-items-center ">
                        <div class="card py-5 px-3">
                            <h5 class="m-0 padding-bottom-40">Enter the verification code you received by SMS.</h5>
                            <span class="mobile-text">We need to verify your phone number before you can publish your page.</span>
                            <p>We sent an SMS to <b class="text-danger">{this.state.phone}</b>. Enter it below and click verify.</p>
                            {
                                this.state.showPhone ?
                                <div class="form-group">
                                    <label>Phone number</label>
                                    <input type="tel" class="form-control" placeholder="Phone number e.g +263783XXXXXX" value={this.state.phone} onChange={(evt) => { this.setState({ phone: evt.target.value }) }} />
                                    <button class="btn btn-block btn-info" onClick={this.updatePhone}><i class="fa fa-save"></i> Save and send code</button>
                                </div> : <></>
                            }
                            <div class="d-flex flex-row mt-5" data-children-count="4">
                                <input type="number" pattern="\d+{6}" class="form-control" autofocus="" value={this.state.otp} onChange={(evt) => {this.setState({ otp: evt.target.value }) }}  placeholder="6 digit OTP sent via SMS" />
                            </div>
                            <button class="btn btn-block btn-primary" onClick={this.verifyPhone}><i class="fa fa-check"></i> Verify code</button>
                            <div class="text-center mt-5 padding-top-10">
                                <span class="d-block mobile-text">Don't receive the code?</span>
                                <a onClick={this.resendCode} href="#" ><span class="font-weight-bold text-danger cursor">Resend</span></a> or 
                                <span class="font-weight-bold text-info cursor"><a href="#" onClick={_ => this.setState({showPhone: true, enableVerify: false})} >Change phone number</a></span>
                            </div>
                        </div>
                    </div>
                </div>
            </>
        )
    }
} 

class VerifyEmailAction extends Component {
    constructor(props) {
        super(props)

        this.state = { verified: this.props.user.emailVerified, showemail: false, enableVerify: true, email: this.props.user.email, otp: '' }
        this.resendCode = () => {
            v1.security.resendEmailCode().then(resp => {
                if (resp.status === 'sent') {
                    alert('Resent the verification code.')
                } else {
                    alert(resp.error)
                }
            })
        }

        this.verifyEmail = () => {
            v1.security.verifyEmail(this.state.otp).then(resp => {
                if (resp.status === 'verified') {

                    window.location.reload()
                } else {
                    alert(resp.error)
                }
            }).catch(err => {
                alert(err.error)
            })
        }

        this.updateEmail = () => {
            v1.security.updateEmail(this.state.email).then(resp => {
                if (resp.status === 'updated') {
                    this.setState({ enableVerify: true, showemail: false })
                } else {
                    alert(resp.error)
                }
            }).catch(err => {
                alert(err.error)
            })
        }
    }

    render() {
        return (
            <>
                <div class="box">
                    <div class="box-header">
                        <h4 class="box-title">Please verify your email address to access this feature.</h4>
                    </div>
                </div>
                <div class="row justify-content-center padding-top-50">
                    <div class="col-md-6 col-xs-12 d-flex justify-content-center align-items-center ">
                        <div class="card py-5 px-3">
                            <h5 class="m-0 padding-bottom-40">Enter the verification code you received by email.</h5>
                            <span class="mobile-text">We need to verify your email address before you are authorized access to this resource.</span>
                            <p>We sent an email to <b class="text-danger">{this.state.email}</b>. Enter it below and click verify.</p>
                            {
                                this.state.showemail ?
                                    <div class="form-group">
                                        <label>* Email address</label>
                                        <input type="email" class="form-control" placeholder="myemail@gmail.com" value={this.state.email} onChange={(evt) => { this.setState({ email: evt.target.value }) }} />
                                        <button class="btn btn-block btn-info" onClick={this.updateEmail}><i class="fa fa-save"></i> Save and send code</button>
                                    </div> : <></>
                            }
                            <div class="d-flex flex-row mt-5" data-children-count="4">
                                <input type="number" pattern="\d+{6}" class="form-control" autofocus="" value={this.state.otp} onChange={(evt) => { this.setState({ otp: evt.target.value }) }} placeholder="6 digit OTP sent to your email" />
                            </div>
                            <button class="btn btn-block btn-primary" onClick={this.verifyEmail}><i class="fa fa-check"></i> Verify code</button>
                            <div class="text-center mt-5 padding-top-10">
                                <span class="d-block mobile-text">Didn't receive the email?</span>
                                <a onClick={this.resendCode} href="#" ><span class="font-weight-bold text-danger cursor">Resend</span></a> or
                                <span class="font-weight-bold text-info cursor"><a href="#" onClick={_ => this.setState({ showemail: true, enableVerify: false })} >Change email address</a></span>
                            </div>
                        </div>
                    </div>
                </div>
            </>
        )
    }
} 

function ChangeUsernameAction(props){
    return (<></>)
}


function VerifyIdentityAlertAction(props) {
    return (
        <div class="box">
            <div class="box-body">
                <div class="row">
                    <div class="col-md-12">
                        <p>
                            <span class="font-weight-bold">Verify you account to access all features offered.</span>
                        </p>
                        <Link to="/creator/verify-identity" class="btn btn-block btn-primary" onClick={props.verifyIdentity}><i class="fa fa-check"></i> Verify identity</Link>
                    </div>
                </div>
            </div>
        </div>
    )
}

class ActionCenter extends Component {

    constructor(props){
        super(props)
    }

    render() {
        let component = <></>
        switch( this.props.action ){
            case 'phone':
                component = <VerifyPhoneNumberAction user={this.props.user} />
            break;
            case 'identity':
                component = <VerifyPhoneNumberAction user={this.props.user} />
                break;
            case 'verify_identity_alert':
                component = <VerifyIdentityAlertAction user={this.props.user}/>
                break;
            case 'username':
                component = <ChangeUsernameAction user={this.props.user} />
                break;
            case 'email':
                component = <VerifyEmailAction user={this.props.user} />
                break;
            default:
                component = <p>Unknown action: {this.props.action} </p>
                break;
        }
        return (
            <>
            {component}
            </>
        );
    }
}

export default ActionCenter;