import React, {Component} from 'react'
import Linkify from 'react-linkify/dist/components/Linkify';
import { Link, Redirect } from 'react-router-dom';
import v1 from "../../api/v1";
import money from "./Amount"
import LoginPage from '../LoginPage';
import Pusher from "pusher-js"

class PaymentUI extends Component{
    constructor(props) {
        super(props)
        var startStep = 'form', formVTitle = 'Pay ' + money.formatUSDToZWL(this.props.amount) + ' to get instant access'
        switch (this.props.purpose) {
            case 'subscribe':
                startStep = 'prompt'
                break;
            case 'support':
                formVTitle = 'Support @' + this.props.creator.username + ' with ' + money.formatUSDToZWL(this.props.amount)
                break;
            default:
                startStep = 'form'
                break;    
        }
        this.state = {
            phone: this.props.user ? this.props.user.phone_number : '',
            email: this.props.user ? this.props.user.email : '',
            fullname: this.props.user ? this.props.user.fullname : 'Anonymous',
            ecocash: this.props.ecocash ? this.props.ecocash : '',
            message: '',
            items: this.props.items ? this.props.items : 1,
            itemName: this.props.itemName ? this.props.itemName : '', 
            purpose: this.props.purpose ? this.props.purpose : 'support',
            nonce: parseInt(Math.random() * 1000),
            password: '',
            amount: this.props.amount,
            payment: {
                _id: 'nil',
                status: 'sending payment',
            },
            form: {},
            step: startStep,
            paid: false,
            payBtnText: '',
            payBtnDisabled: false,
            payBtnLoading: false,
            payBtnClass: 'btn-primary',
            payBtnIcon: 'phone',
            formViewTitle: formVTitle,
            formViewMessage: '',
            pollCancelFunc: () => { },
            showPayBtn: false,
            showCancelBtn: false,
            error: '',
            pollCounter: 0,
            pollInterval: 15,
            polling: false,
            signature: '',
            ts: 0,
            pusher: null,
            cancelBtnText: 'Close this window',
            cancelBtnIcon: 'times',
            gateway: 'ecocash'
        }

        

        //subscribe prompt
        this.subscribePromptView = this.subscribePromptView.bind(this)

        //forms
        this.formView = this.formView.bind(this)
        this.subscribeFormView = this.subscribeFormView.bind(this)
        this.supportFormView = this.supportFormView.bind(this)
        this.payPerViewFormView = this.payPerViewFormView.bind(this)
        this.payForServiceFormView = this.payForServiceFormView.bind(this)

        //validations
        this.validatePayPerViewForm = this.validatePayPerViewForm.bind(this)
        this.validateServiceForm = this.validateServiceForm.bind(this)
        this.validateSubscribeForm = this.validateSubscribeForm.bind(this)
        this.validateSupportForm = this.validateSupportForm.bind(this)

        this.errorView = this.errorView.bind(this)
        this.processingView = this.processingView.bind(this)
        this.successView = this.successView.bind(this)

        this.submitPayForm = this.submitPayForm.bind(this)
        this.cancelTransaction = this.cancelTransaction.bind(this)
        this.payBtnAction = this.payBtnAction.bind(this)
        this.showLoginPage = this.showLoginPage.bind(this)

        this.pollPayment = this.pollPayment.bind(this)
        this.onPaid = this.onPaid.bind(this)
        this.onCancelled = this.onCancelled.bind(this)
        this.onSent = this.onSent.bind(this)
        this.triggerOtherPayMethod = this.triggerOtherPayMethod.bind(this)
        this.formCreateAccount = this.formCreateAccount.bind(this)

        console.log(this.state.purpose)
    }

    componentDidCatch(error, info) {
        console.log(error, info)
    }

    pollPayment() {

        this.setState({ pollCounter: this.state.pollCounter + 1 })
        if (this.state.payment.status === 'cancelled') {
            this.onCancelled(this.state.payment)
            clearInterval(this.state.pollCancelFunc)
            return
        }

        if (this.state.payment.status === 'paid') {
            this.onPaid(this.state.payment)
            clearInterval(this.state.pollCancelFunc)
            return
        }

        if (this.state.pollCounter <= this.state.pollInterval && !this.state.polling) {
            this.setState({
                payBtnText: 'Checking payment in '+ (this.state.pollInterval-this.state.pollCounter) + ' seconds',
            })
            return
        } else {
            this.setState({
                payBtnText: 'Checking payment ...',
                pollCounter: 0, polling: true, 
            })
        }
        
        v1.payments.pollPayment(this.state.payment._id, this.state.signature, this.state.ts).then(resp => {
            this.setState({ polling: false, payment: resp})
            if (resp.error) {
                this.setState({
                    step: 'error',
                    error: resp.error,
                    payBtnText: 'Retry this transaction now',
                    payBtnDisabled: false,
                    payBtnLoading: false,
                    payBtnClass: 'btn-primary',
                    payBtnIcon: 'phone',
                    showCancelBtn: true,
                })
                return
            }

            if (resp.status === this.state.payment.status) {
                return
            }
            switch (resp.status) {
                case 'paid': this.onPaid(resp); break;
                case 'cancelled': this.onCancelled(resp); break;
                case 'sent': this.onSent(resp);  break;
                default: console.log('Unknown transaction service response: ', resp)
            }
        }).catch(err => {
            this.setState({polling: false, error: err})
        })
        
    }

    onPaid(resp) {
        console.log(resp)
        console.log(this.state)
        if (this.state.paid) {
            return
        }
        this.setState({paid: true, payment: resp, step: 'success', showCancelBtn: false, cancelBtnText: '', payBtnIcon: "check", payBtnText: 'Payment successful!', payBtnDisabled: true, payBtnLoading: false, payBtnClass: 'btn-success' })
        if (this.state.purpose === 'pay_per_view') {
            setTimeout(_ => window.location = resp.unlock_code, 1000)
            alert({
                        toast: true,
                        icon: 'success',
                        timer: 5000,
                        title: 'Content Unlocked',
                        text: 'Your payment has been received, you can now view this content'
                    })
        }
        if (this.state.purpose === 'subscribe') {
            alert({
                toast: true,
                icon: 'success',
                timer: 5000,
                title: 'Subscription active :)',
                text: 'Your payment has been received, you can now view this content'
            })
            this.setState({step: 'goto_content'})
        }
        if (this.state.purpose === 'support') {
            this.setState({step: 'show_thankyou'})   
        }
        clearInterval(this.state.pollCancelFunc)
    
    }

    onCancelled(resp) {
        this.setState({ payment: resp, step: 'error', payBtnIcon: "times", error: 'Payment cancelled', payBtnText: 'Retry this transaction ', payBtnDisabled: false, payBtnLoading: false, payBtnClass: 'btn-primary' })
        this.state.pusher.disconnect()
        clearInterval(this.state.pollCancelFunc)
    }

    onSent(resp) {
        this.setState({ payment: resp, showCancelBtn: true, step: 'processing', payBtnIcon: "history", payBtnText: 'Payment sent!', payBtnDisabled: true, payBtnLoading: true, payBtnClass: 'btn-default' })
    }

    payBtnAction(evt) {
        evt.preventDefault()
        if (this.state.step === 'prompt') {
            //add 3 months to expiry date
            let expiryDate = new Date()
            expiryDate.setMonth(expiryDate.getMonth() + 3)
            let msg = 'Subscribe now for only '+ money.formatUSD(this.props.amount) + ' for 3 months'
            this.setState({
                step: 'form',
                formViewTitle: this.props.user.logged_in ? 'Notifications will be sent to ' + this.state.email : 'Create your account to subscribe.',
                formViewMessage: this.props.user.logged_in ? msg : <p>{msg}<br/>Your account will be created when your payment is received.</p>,
                showPayBtn: false,
                showCancelBtn: false,
            })
            return
        }
        if (this.state.step === 'form') {
            this.submitPayForm()
            return
        } else if (this.state.step === 'processing') {
            this.cancelTransaction()
            return
        } else if (this.state.step === 'error') {
            this.setState({step: 'form', showPayBtn: false, showCancelBtn: false, nonce: parseInt(Math.random() * 1000)})
            return
        }
    }


    triggerOtherPayMethod() {
        alert('Sorry. Only Ecocash payments are active at the moment.')
    }

    showLoginPage() {
        this.setState({
            step: 'login',
            showPayBtn: false,
            showCancelBtn: false,
        })
    }

    submitPayForm() {
        let b = false
        
        switch (this.props.purpose) {
            case 'subscribe':
                b = this.validateSubscribeForm()
                break;
            case 'support':
                b = this.validateSupportForm()
                break;
            case 'pay_per_view':
                b = this.validatePayPerViewForm()
                break;
            case 'service':
                b = this.validateServiceForm()
                break;
            default:
                b = false
                break;    
        }
        if (!b) {
            return
        }
        this.setState({
            pusher: new Pusher('d3c8fc4f4d1b77ff0011', {
                cluster: 'eu'
            })
        });

        this.setState({
            payBtnText: 'Sending payment...',
            payBtnDisabled: true,
            step: 'processing',
            payBtnLoading: true,
            payBtnClass: 'btn-warning',
            payBtnIcon: 'spin fa-spinner',
            showCancelBtn: true,
            showPayBtn: true,
        })

        let payment = {}
        
        switch (this.state.purpose) {
            case 'subscribe':
                payment = {
            email: this.state.email,
            fullname: this.state.fullname,
            items: this.props.items ? this.props.items : 1,
            phone: this.state.ecocash,
            notification_phone: this.state.phone,
            password: this.state.password,
            campaign: this.props.content && this.props.content._id ? this.props.content._id : undefined,
            supporter: this.props.user._id,
            creator: this.props.creator.username,
            gateway: this.state.gateway,
                }
                break;
            case 'support':
                payment = {
                    amount: this.state.amount,
                    email: this.state.email,
                    fullname: this.state.fullname,
                    message: this.state.message,
                    phone: this.state.ecocash,
                    notification_phone: this.state.phone,
                    items: this.state.items,
                    campaign: this.props.content && this.props.content._id ? this.props.content._id : undefined,
                    creator: this.props.creator.username,
                    gateway: this.state.gateway,
                }
                break;
            case 'pay_per_view':
                payment = {
                    email: this.state.email,
                    fullname: this.state.fullname,
                    phone: this.state.ecocash,
                    notification_phone: this.state.phone,
                    campaign: this.props.content._id,
                    supporter: this.props.user._id,
                    creator: this.props.creator.username,
                    gateway: this.state.gateway,
                }
                break;
            case 'service':
                payment = {
                    ...this.props.form,
                    ecocash: this.state.ecocash,
                    service: this.props.service._id,
                    supporter: this.props.user._id,
                    creator: this.props.creator.username,
                    gateway: this.state.gateway,
                }
                break;
            default:
                alert('Unsupported action!')
                window.location.href = "/logout"
                break;    
        }
        

        v1.payments.dispatchPayment(this.state.purpose, payment, this.state.nonce).then(resp => {
            if (resp.error) {
                this.setState({
                    step: 'error',
                    error: resp.error,
                    payBtnText: 'Retry this transaction ' ,
                    payBtnDisabled: false,
                    payBtnLoading: false,
                    payBtnClass: 'btn-primary',
                    payBtnIcon: 'phone',
                    showCancelBtn: true,
                    cancelBtnIcon: 'times',
                    cancelBtnText: 'Close this window',
                })
                return
            } else {
                //also do http polling every 10seconds
                let pcf = setInterval(_ => this.pollPayment(resp._id, resp.ts, resp.signature, 0), 1000)
                this.setState({
                    pollCancelFunc: pcf,
                    step: 'processing',
                    payment: resp.payment,
                    signature: resp.signature,
                    ts: resp.ts,
                    payBtnText: 'Payment sent!',
                    cancelBtnText: 'Cancel this transaction',
                    payBtnDisabled: true,
                    payBtnLoading: true,
                    payBtnClass: 'btn-default',
                })

                try {
                    let sub = this.state.pusher.subscribe(resp._id)
                    sub.bind('sent', this.onSent)
                    sub.bind('paid', this.onPaid)
                    sub.bind('cancelled', this.onCancelled)
                } catch (e) {
                    console.log("Failed to register pusher: " + JSON.stringify(e))
                }
                
                
                return
            }
        }).catch(err => {
            this.setState({
                    step: 'error',
                    error: err.error ? err.error : err,
                    payBtnText: ':( Transaction failed',
                    payBtnDisabled: true,
                    payBtnLoading: false,
                    payBtnClass: 'btn-primary',
                    payBtnIcon: 'phone',
                    showCancelBtn: true,
                    cancelBtnIcon: 'times',
                    cancelBtnText: 'Close this window'
                })
        })
    }

    cancelTransaction() {
        if (this.state.step === 'processing' || this.state.payment.status === 'sent' || this.state.payment.status === 'queued') {
            alert('To cancel this transaction, simply ignoring the message on the Ecocash number you specified. This page will update after a few seconds.')
            return
        }
        this.props.onClose()
    }

    formCreateAccount() {
        return(
<div class="row row-1 ">
    <div class="col-2"><i class="fa fa-user"></i></div>
    <div class="col-10"><input type="text" placeholder="Your name." class="card-subpay-input" value={this.state.fullname} onChange={evt => this.setState({fullname: evt.target.value})}/></div>
    <div class="col-2 mt-10"><i class="fa fa-envelope"></i></div>
    <div class="col-10 mt-10"><input type="email" placeholder="Your email address" class="card-subpay-input" value={this.state.email} onChange={evt => this.setState({email: evt.target.value})}/></div>
    <div class="col-2 mt-10"><i class="fa fa-lock"></i></div>
    <div class="col-10 mt-10"><input type="password" placeholder="Enter account password" class="card-subpay-input" value={this.state.password} onChange={evt => this.setState({password: evt.target.value})} /></div>
</div>
        )
    }

    validateSubscribeForm() {
        let phone = this.state.ecocash
        if (phone.startsWith("0")) {
            phone = phone.substring(1);
        }
        if (phone.length < 9) {
            alert("Please enter a valid phone number")
            return false
        }
        if (!phone.startsWith("77") && !phone.startsWith("78")) {
            alert("Phone number number start with 77 or 78")
            return false
        }

        if (this.state.fullname === '') {
            this.setState({ fullname: 'Anonymous' })
        }

        if (!v1.util.validateEmail(this.state.email)) {
            alert("Please enter a valid email address")
            return false
        }
        return true
    }

    formView() {
        return (
<form onSubmit={evt => evt.preventDefault()} >
    <span class="card-subpay-header">{this.state.formViewTitle}</span>
<>
    {this.state.purpose === 'subscribe' ? this.subscribeFormView() : <></>}
    {this.state.purpose === 'pay_per_view' ? this.payPerViewFormView() : <></>}
    {this.state.purpose === 'support' ? this.supportFormView() : <></>}
    {this.state.purpose === 'service' ? this.payForServiceFormView() : <></>}
</>
    <div class="row row-1">
        <div class="col-2"><img alt="" class="img-fluid" src="/assets/img/ecocash.svg" /></div>
        <div class="col-7"><input type="tel" placeholder="Your ecocash phone number" class="card-subpay-input" value={this.state.ecocash} maxLength={10} onChange={evt => this.setState({ecocash: evt.target.value})} /></div>
        <div class="col-3 d-flex justify-content-center"><button class="btn btn-danger"  onClick={this.payBtnAction}> <i class="fa fa-mobile"></i> Pay</button></div>
    </div>
    <span class="card-subpay-header">- Or -</span>
    <button class="btn btn-block btn-default mt-5 mb-10" onClick={this.triggerOtherPayMethod}><i class="fa fa-credit-card"></i> Pay with Mastercard/VISA</button>
    {this.state.purpose === 'subscribe' ? <div class="row alert alert-info row-1">
        <p>{this.state.formViewMessage}</p>
    </div> : <></>}
</form>
        )
    }

    subscribeFormView() {
        let component = this.props.user.logged_in ? this.payPerViewFormView() : this.formCreateAccount()
        return component
    }

    validatePayPerViewForm() {
        return this.validateSubscribeForm()
    }

    payPerViewFormView() {
        return (
<div class="row row-1 ">
    <div class="col-2"><i class="fa fa-user"></i></div>
    <div class="col-10"><input type="text" placeholder="Your name." class="card-subpay-input" value={this.state.fullname} onChange={evt => this.setState({fullname: evt.target.value})}/></div>
    <div class="col-2 mt-10"><i class="fa fa-envelope"></i></div>
    <div class="col-10 mt-10"><input type="email" placeholder="Your email address" class="card-subpay-input" value={this.state.email} onChange={evt => this.setState({email: evt.target.value})}/></div>
</div>    
)
    }


    validateSupportForm() {
        let b = this.validateSubscribeForm()
        if(b){
            if (this.state.message.length > 1024) {
                alert("Message is too long")
                return false
            }
        }
        return b
    }

    supportFormView() {
        return (
<div class="row row-1 ">
    <div class="col-2"><i class="fa fa-user"></i></div>
    <div class="col-10"><input type="text" placeholder="Your name." class="card-subpay-input" value={this.state.fullname} onChange={evt => this.setState({fullname: evt.target.value})}/></div>
    <div class="col-2 mt-10"><i class="fa fa-envelope"></i></div>
    <div class="col-10 mt-10"><input type="email" placeholder="Your email address" class="card-subpay-input" value={this.state.email} onChange={evt => this.setState({email: evt.target.value})}/></div>
    <div class="col-2 mt-10"><i class="fa fa-heart" style={{color: "pink"}}></i></div>
    <div class="col-10 mt-10"><textarea placeholder="Optional leave a message. Your message can be seen publicly" class="card-subpay-input" rows={5} style={{height:'auto'}} value={this.state.message} onChange={evt => this.setState({message: evt.target.value})} ></textarea></div>
</div>
        )
    }

    validateServiceForm() {
        return true
    }

    payForServiceFormView() {
        return (<>
        <table class="table table-responsive table-striped">
            <tbody>
                <tr>
                    <td>Order for </td>
                    <td><Link to={""} target="_blank">{this.props.service.title}</Link></td>

                </tr>
                
                <tr>
                    <td>Description</td>
                    <td>{this.props.service.description}</td>
                </tr>
                <tr>
                    <td>Fullname</td>
                    <td>{this.props.form.fullname}</td>
                </tr>
                <tr>
                    <td>Email</td>
                    <td>{this.props.form.email}</td>
                </tr>
                <tr>
                    <td>Phone number</td>
                    <td>{this.props.form.phone}</td>
                </tr>
                <tr>
                    <td>{this.props.form.question}</td>
                    <td>{this.props.form.answer}</td>
                    </tr>
                <tr>
                    <td>Items left </td>
                    <td>{this.props.service.service.quantity_available}</td>

                </tr>
            </tbody>
        </table>
        </>)
    }

    subscribePromptView() {
        return (
<form class="">
   <div class="card">
      <div class="row justify-content-center">
         <div class="">
            <div class="mb-30">
               <div class="">
                  <img src={this.props.creator.subscriptions.url} class="img thumb" alt="headline img" />
                  <p>Access this and other exclusive content by @{this.props.creator.username} by subscribing to their account.</p>
                  <p></p>
               </div>
            </div>
            <div class="form-group row justify-content-center mb-0">
               <div class="col-md-12 px-3">
                   <button class="btn btn-block btn-black btn-block rm-border" onClick={this.payBtnAction}>
                       <i class="fa fa-credit-card"></i>Subscribe for {money.formatUSD(this.props.amount)}
                    </button>
                </div>
                </div>
                {this.props.user.logged_in ? <></> :
                    <div class="form-group row justify-content-center mb-0">
                        <div class="col-md-12 px-3 mt-1">
                            <Link to={"/auth/login?redirect=" + this.props.target} class="btn btn-default btn-block">
                                <p class="loginbtn-clear">Already subscribed? Login to your account</p>
                            </Link>
                        </div>
                    </div>
                }
         </div>
      </div>
   </div>
</form>
            )
    }

    processingView() {
        return (
<div class="row">
    <div class="col-12">
        <h4>Processing payment - {this.state.payment  ?  this.state.payment.status : 'pending'}</h4>
        <div class="row justify-content-center">
            <div class="col-md-6">
                <img src="/assets/img/loading.gif" alt="loading" class="img-responsive" />
            </div>
        </div>
        <strong>Keep your eyes on your phone! A authorization prompt should appear for you to pay</strong>
        <p>
            <small>Payment ID: <b>{this.state.payment  ? this.state.payment._id : '--- ---- ---- ----'}</b></small>
        </p>
        <p>
            <small>Ecocash Number: <b>{this.state.ecocash ? this.state.ecocash : '07XXYYYZZZ'}</b></small>
        </p>
        <p>
            <small>Email: <b>{this.state.email ? this.state.email : '________@____.__'}</b></small>
        </p>
    </div>
</div>
        )
    }

    successView() {
        return (
        <div class="row">
            <div class="col-md-12">
                <div class="row justify-content-center">
                    <div class="col-6">
                        <img src="/assets/img/ecocash_paid.gif" class="img-responsive" alt="success" />
                    </div>
                </div>
                <h4>Payment Received</h4>
                <ul class="timeline">
                    <li>
                            <h6 style={{ "paddingLeft": "20px", "paddingTop": "20px" }}>Thank you message from @{this.props.creator.username}</h6>
                            {this.props.items ?
                                <p style={{ "paddingLeft": "20px" }} class="float-">you bought me {this.props.items} {this.props.itemName}(s)</p> : <></>
                            }
                        <p  style={{ "paddingLeft": "20px", "paddingTop": "20px", whiteSpace: "pre-wrap" }}>
                            <Linkify
                                componentDecorator={(decoratedHref, decoratedText, key) => (
                                    <a target="blank" style={{color: 'red', fontWeight: 'bold'}} rel="noopener" target="_blank" href={decoratedHref} key={key}>
                                        {decoratedText}
                                    </a>
                                )}
                            >{ this.state.payment.thank_you }</Linkify>    
                        </p>
                    </li>
                </ul>    
            </div>
                        
        </div>
        )
    }


    errorView(){
        return (
<div class="row">
    <div class="col-12">
        <h4>Payment failed :(</h4>
        <div class="row justify-content-center">
            <div class="col-md-6">
                <img src="/assets/img/payment_failed.gif" alt="failed" class="img-responsive" />
            </div>
        </div>
        <strong>{JSON.stringify(this.state.error)}</strong>
    </div>
</div>
        )
    }

    render() {
        return (<div class="row">
            <div class="col-12">
                <div class="tab-content">
                    <div class="tab-pane fade show active" id="nav-tab-card">
                        {this.state.step === 'prompt' ? this.subscribePromptView() : <></>}
                        {this.state.step === 'goto_content' ? <Redirect to={this.props.target} /> : <></>}
                        {this.state.step === 'login' ? <LoginPage onLogin={_ => this.setState({ step: 'goto_content' })} /> : <></>}
                        
                        {this.state.step === 'form' ? this.formView() : <></>}

                        {this.state.step === 'show_thankyou' ? this.successView() : <></>}

                        {this.state.step === 'processing' ? this.processingView() : <></>}
                        {this.state.step === 'success' ? this.successView() : <></>}
                        {this.state.step === 'error' ? this.errorView() : <></>}
            
          

                    </div>
                    {this.state.showPayBtn ? <button onClick={this.payBtnAction} class={"btn btn-block " + this.state.payBtnClass} disabled={this.state.payBtnDisabled}  >
                        <i class={"fa fa-" + this.state.payBtnIcon}></i> {this.state.payBtnText}
                    </button> : <></>
                    }
                    {this.state.showCancelBtn ?
                        <button onClick={this.cancelTransaction} class="btn btn-block btn-sm btn-default margin-top-10" >
                            <i class={"fa fa-" + this.state.cancelBtnIcon}></i> {this.state.cancelBtnText}
                        </button>
                    : <></>}
                            
                </div>
            </div>
        </div>)
    }
}


export default PaymentUI;