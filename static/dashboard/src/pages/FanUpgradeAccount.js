import React, { Component } from "react";
import { Link } from "react-router-dom"


import v1 from "../api/v1";
import InlineUpload from "../components/InlineUploader";

class FanUpgradeAccount extends Component {
    constructor(props) {
        v1.page.set({ title: 'Verify your identity for more benefits' })
        v1.page.event('Verify', 'Identity', 'Page')
        super(props)

        this.state = {
            payVerifyMsg: 'Pay USD $1.00 now',
            payVerifyDisabled: false,
            payVerifyPhone: '',
            payNextPollSecond: 0
        }

        this.registerPaymentChecker = this.registerPaymentChecker.bind(this)
        this.verifyByPay = this.verifyByPay.bind(this)

    }

    registerPaymentChecker(res, callback, interval = 1000, pollInterval = 15000, counter = 0) {
        if (counter >= pollInterval) {
            counter = 0
            v1.payments.pollPayment(res._id, res.signature, res.ts).then(payment => {
                if (!callback(payment)) {
                    counter += interval
                    setTimeout(() => this.registerPaymentChecker(res, callback, 1000, 15000, counter), 1000)
                }
            }).catch(err => {
               //we move on 
                console.log(err)
                setTimeout(() => this.registerPaymentChecker(res, callback, 1000, 15000, counter), 1000)
            })
        } else {
            this.setState({ payVerifyMsg: 'Checking payment in ' + ((pollInterval - (counter + interval)) / 1000) + ' seconds' })
            setTimeout(() => this.registerPaymentChecker(res, callback, interval, pollInterval, counter + interval), interval)
        }
    }

    verifyByPay() {
        let phone = this.state.payVerifyPhone ? this.state.payVerifyPhone : '';
        if (phone.startsWith("0")) {
            phone = phone.substring(1);
        }
        if (phone.length < 9) {
            alert("Please enter a valid phone number")
            return
        }
        if (!phone.startsWith("77") && !phone.startsWith("78")) {
            alert("Phone number number start with 77 or 78")
            return
        }
        this.setState({ payVerifyDisabled: true, payVerifyMsg: 'Initiating payment...' })

        v1.security.verifyByPayment(phone).then(res => {
            if (res.error) {
                alert('Try again. Failed with error: ' + res.error)
                this.setState({ payVerifyDisabled: false, payVerifyMsg: 'Pay USD $1.00 now' })
                return
            }
            if (res.status === 'queued') {
                this.setState({ payVerifyDisabled: true, payVerifyMsg: 'Payment pending...' })
                alert("Payment initiated. Please keep your phone unlocked and wait for prompt to arrive.\n\nThis page will automatically reload once payment has been received.")
                this.registerPaymentChecker(res, (payment) => {
                    if (payment.status === 'paid') {
                        this.setState({ payVerifyDisabled: true, payVerifyMsg: 'Payment Received! - Account verified !' })
                        alert('Payment received. You account has been verified.')
                        setTimeout(_ => { window.location = '/creator/dashboard' }, 1000)
                        return true
                    } else if (payment.status === 'cancelled') {
                        this.setState({ payVerifyDisabled: false, payVerifyMsg: 'Pay USD $1.00 now' })
                        alert('Payment failed. Please try again or use a different verification method.')
                        return true
                    } else if (payment.status === 'sent' || payment.status === 'queued') {
                        this.setState({ payVerifyDisabled: true, payVerifyMsg: 'Payment pending...' })
                        return false
                    }
                })
            } else {
                alert('Received unknown status: ' + res.status)
                this.setState({ payVerifyDisabled: false, payVerifyMsg: 'Pay USD $1.00 now' })
            }
            
        }).catch(err => {
            alert(err.error)
            this.setState({ payVerifyDisabled: false, payVerifyMsg: 'Pay USD $1.00 now' })
        })
        
    }

    render() {
        return (
            <>
            <div class="box">
                <div class="box-header">
                    <h6 class="box-title">Choose one of three methods to verify your account.</h6>
                </div>
            </div>
            <div class="box">
                <div class="box-header">
                    <h6 class="box-title">Option 1. Verify by one time payment.</h6>
                </div>
                <div class="box-body">
                    <div class="row">
                            <div class="col-md-12">
                                <div class="form-group">
                                    <label for="pwd">* Ecocash Phone number</label>
                                    <div class="input-group ">
                                        <span class="input-group-text border-0" id="basic-addon3">+263</span>
                                        <input onChange={(evt) => {this.setState({payVerifyPhone : evt.target.value }) }} value={this.state.payVerifyPhone} type="text" readOnly={this.state.payVerifyDisabled} maxLength={10} class="form-control rounded" id="basic-url" placeholder="77/78 Ecocash phone number" />
                                    </div>
                                    <br />
                                    <p>
                                        To verify your account, you will need to pay a one time fee of USD $1.00(ZWL equiv).<br />
                                        Verification is instant.
                                    </p>
                                    <button class="btn btn-block btn-primary" onClick={this.verifyByPay} disabled={this.state.payVerifyDisabled}><i class="fa fa-credit-card"></i>{ this.state.payVerifyMsg }</button>
                                </div>
                            </div>
                        </div>
                </div>
                </div>
                <div class="box">
                <div class="box-header">
                    <h6 class="box-title">Option 2. Enter invite code to verify account</h6>
                </div>
                <div class="box-body">
                    <div class="row">
                            <div class="col-md-12">
                                <div class="form-group">
                                    <label for="pwd">* 8 Digit invite code</label>
                                    <div class="input-group ">
                                        <input type="text" maxLength={8} class="form-control rounded" id="basic-url" placeholder="8 digit invite code" />
                                    </div>
                                    <p>
                                        Instantly verify your account by entering an invite code you received.<br />
                                        The invite code will only work on your account.<br />
                                        Verification is instant
                                    </p>
                                    <button class="btn btn-block btn-success" ><i class="fa fa-key"></i> Verify invite code </button>
                                </div>
                            </div>
                        </div>
                </div>
                </div>
                <div class="box">
                <div class="box-header">
                    <h6 class="box-title">Option 3. Upload your ID Document</h6>
                </div>
                <div class="box-body">
                    <div class="row">
                            <div class="col-md-12">
                                <div class="form-group">
                                    <InlineUpload type="image" maxNumberOfFiles={1} purpose="verify_identity" onUploaded={this.onUploaded} allowedTypes={['image/*']} />
                                    <p>
                                        Valid identity documents are: Passport, National ID, Drivers License, Bank statement.
                                        <ul>
                                            <li>Verification takes up to 72 hours.</li>
                                            <li>We might contact you directly for additional verification.</li>
                                            <li>Once verified we securely store your documents in cold storage(device not connected to the internet)</li>
                                        </ul>
                                    </p>
                                    <button class="btn btn-block btn-danger" ><i class="fa fa-user"></i> Submit verification request. </button>
                                </div>
                            </div>
                        </div>
                </div>
                </div>
                </>
           
        );
    }
}

export default FanUpgradeAccount;

