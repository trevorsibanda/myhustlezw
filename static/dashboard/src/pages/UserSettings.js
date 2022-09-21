import { Component } from "react";


import {Link, NavLink, Switch, Route} from "react-router-dom"
import v1 from "../api/v1";
import Preloader from "../components/PreLoader";
import RichEditor from "../components/RichEditor"
import PageCoverEditor from "../components/PageCoverEditor"
import PageSettings from "./PageSettings";
import Wallet from "./Wallet";

import InlineUpload from "../components/InlineUploader";
import Restrict from "../components/Restrict";
import PriceButtonGroup from "../components/payments/PriceButtonGroup";


class SecuritySettings extends Component {
    constructor(props){
        super(props)

        this.state = {
            current: '',
            new: '',
            new2: '',
        }
        this.updatePassword = this.updatePassword.bind(this)
    }
    
    updatePassword(){
        if(this.state.new.length < 6 ){
            return alert('Password should be atleast 6 characters long.')
        }
        if(this.state.new != this.state.new2){
            return alert('Passwords do not match!')
        }
        v1.security.updatePassword(this.state.current, this.state.new).then(resp => {
            if(resp.status === 'changed'){
                alert('Successfully changed password. This page will close in 5 seconds and you will have to log in again')
                setTimeout(_ => window.location = "/logout?change_password", 3000)
            }else {
                alert(resp.error)
            }
        }).catch(err => {
            alert(err.error)
        })
    }

    render() {
        return (
            <div class="box-body">
                <h3>Security settings</h3>
                <p></p>
                <div class="form-group">
                    <label>Current Password</label>
                    <br />
                    <small>Do not enter this, unless you plan on changing your password</small>
                    <div class="input-group mb-3">
                        <div class="input-group-prepend">
                            <span class="input-group-text"><i class="fa fa-lock"></i></span>
                        </div>
                        <input type="password" class="form-control" placeholder="Confirm Password" value={this.state.current} onChange={(evt) => { this.setState({ current: evt.target.value }) }} />
                    </div>
                </div>
                <div class="form-group">
                    <label>New Password</label>
                    <div class="input-group mb-3">
                        <div class="input-group-prepend">
                            <span class="input-group-text"><i class="fa fa-lock"></i></span>
                        </div>
                        <input type="password" class="form-control" placeholder="New Password (minimum of 6 characters)" value={this.state.new} onChange={(evt) => { this.setState({ new: evt.target.value }) }} />
                    </div>
                </div>
                <div class="form-group">
                    <label>Confirm New Password</label>
                    <div class="input-group mb-3">
                        <div class="input-group-prepend">
                            <span class="input-group-text"><i class="fa fa-lock"></i></span>
                        </div>
                        <input type="password" class="form-control" placeholder="Confirm Password" value={this.state.new2} onChange={(evt) => { this.setState({ new2: evt.target.value }) }} />
                    </div>
                    {this.state.new != this.state.new2 ? <p class="text-danger">Your passwords do not match.</p> : <strong class="text-success">Passwords match</strong>}
                </div>
                <hr/>
                <div class="row justify-content-center" >
                    <div class="col-md-4" >
                        
                        <button class="btn btn-block btn-primary" onClick={this.updatePassword}  ><i class="fa fa-warning"></i> Update password</button>
                    </div>
                </div>
            </div>
        )
    }
}

class NotificationSettings extends Component{
    constructor(props){
        super(props)
        this.state = this.props.user.notifications

        this.saveChanges = this.saveChanges.bind(this)
    }

    saveChanges() {
        v1.user.updateNotifications(this.state).then(resp => {
            if (resp.status === 'saved') {
                alert('Successfully saved your changes')
            } else {
                alert(resp.error)
            }
        }).catch(alert)
    }

    render(){
        return (
            <div class="box">
                <div class="box-body">
                    <div class="form-group">
                        <h3>Your notification settings</h3>
                        <label>Notifications</label>
                        <div class="c-inputs-stacked">
                            <input type="checkbox" checked={this.state.subscriber} onChange={(evt) => { this.setState({ subscriber: evt.target.checked}) }} id="checkbox_123" />
                            <label for="checkbox_123" class="block">Notify me of new supporters.</label>
                            <input type="checkbox" checked={this.state.wallet_credit} onChange={(evt) => { this.setState({ wallet_credit: evt.target.checked }) }} id="checkbox_234" />
                            <label for="checkbox_234" class="block">Notify me when I make money.</label>
                            <input type="checkbox" checked={this.state.login} onChange={(evt) => { this.setState({ login: evt.target.checked }) }} id="checkbox_345" />
                            <label for="checkbox_345" class="block">Notify me when I login</label>
                            <input type="checkbox" checked={this.state.service_scheduled} onChange={(evt) => { this.setState({ service_scheduled: evt.target.vchecked}) }} id="checkbox_456" />
                            <label for="checkbox_456" class="block">Notify me when I a new service form is completed</label>
                        </div>
                    </div>

                    <div class="row justify-content-center" >
                        <div class="col-md-4" >
                            <button class="btn btn-block btn-primary" onClick={this.saveChanges} ><i class="fa fa-check"></i> Update notification preferences</button>
                        </div>
                    </div>
                </div>
            </div>
        )
    }
}

class PersonalDetailsSettings extends Component {
    constructor(props) {
        super(props)
        this.state = {
            username: this.props.user.username,
            fullname: this.props.user.fullname,
            email: this.props.user.email,
            phone_number: this.props.user.phone_number,
        }

        this.saveChanges = this.saveChanges.bind(this)
    }

    saveChanges() {
        v1.user.saveBasics(this.state).then(resp =>{
            if(resp.status === 'saved'){
                alert('Successfully saved your changes. Your page will now reload')
                window.location.reload()
            }else{
                alert(resp.error)
            }
        }).catch(alert)
    }

    render(){
        return (
            <div class="box">
                <div class="box-body">
                    <h3>Personal details</h3>
                    <div class="form-group">
                        <label>User Name</label>
                        <div class="input-group mb-3">
                            <div class="input-group-prepend">
                                <span class="input-group-text">https://myhustle.co.zw/@</span>
                            </div>
                            <input type="text" class="form-control" value={this.props.user.username} readonly disabled placeholder="Username" />
                            
                        </div>
                        <p>
                            <small>You can change your username in your <Link to="/creator/settings/mypage" className="text-danger">Page settings</Link></small>
                        </p>
                    </div>
                    <div class="form-group">
                        <label>Fullname</label>
                        <div class="input-group mb-3">
                            <div class="input-group-prepend">
                                <span class="input-group-text"><i class="fa fa-user"></i></span>
                            </div>
                            <input type="text" class="form-control" placeholder="Your full name" value={this.state.fullname} onChange={(evt) => { this.setState({ fullname: evt.target.value }) }} />
                        </div>
                    </div>
                    <div class="form-group">
                        <label>Email address</label>
                        <div class="input-group mb-3">
                            <div class="input-group-prepend">
                                <span class="input-group-text"><i class="fa fa-envelope"></i></span>
                            </div>
                            <input type="email" class="form-control" placeholder="Email" value={this.state.email} onChange={(evt) => { this.setState({ email: evt.target.value }) }} />
                        </div>
                    </div>
                    <div class="form-group">
                        <label>Phone number</label>
                        <div class="input-group mb-3">
                            <div class="input-group-prepend">
                                <span class="input-group-text"><i class="fa fa-phone"></i></span>
                            </div>
                            <input type="tel" class="form-control" placeholder="+263783xxxxxx" value={this.state.phone_number} onChange={(evt) => { this.setState({ phone_number: evt.target.value }) }} />
                        </div>
                    </div>
                    <div class="row justify-content-center" >
                        <div class="col-md-4" >
                            <button class="btn btn-block btn-primary" onClick={this.saveChanges} ><i class="fa fa-check"></i> Save changes</button>
                        </div>
                    </div>
                </div>
            </div>
        )
    }
}

class MyLinks extends Component {

    render() {
        return (
            <>
            <h5>My links</h5>
                <div class="box">
                    <div class="box-header with-border">
                        <h4 class="box-title">My links</h4>
                    </div>
                    <div class="box-body p-0">
                        <div class="media-list media-list-hover media-list-divided">
                            <Link to='/auth/login' class="media media-single">
                                <i class="font-size-18 mr-0 flag-icon flag-icon-us"></i>
                                <span class="title">My page </span>
                                <span class="badge badge-pill badge-secondary">https://myhustle.co.zw/@{this.props.user.username}</span>
                            </Link>

                            <a class="media media-single" href="#">
                                <i class="font-size-18 mr-0 flag-icon flag-icon-ba"></i>
                                <span class="title">Subscribe</span>
                                <span class="badge badge-pill badge-primary">https://myhustle.co.zw/@{this.props.user.username}/subscribe</span>
                            </a>

                            <a class="media media-single" href="#">
                                <i class="font-size-18 mr-0 flag-icon flag-icon-ch"></i>
                                <span class="title">Buy me a {this.props.user.page.donation_item}</span>
                                <span class="badge badge-pill badge-info">https://myhustle.co.zw/@{this.props.user.username}/buymeacoffee</span>
                            </a>

                            
                        </div>
                    </div>
                </div>
            </>)
    }
}

class AdvancedPageSettings extends Component {
    constructor(props){
        super(props)
        this.state = this.props.user.page;
        this.state.subscription_period = this.props.user.subscriptions.period
        this.state.subscription_price = this.props.user.subscriptions.price
        this.state.subscriptions_active = this.props.user.subscriptions.active
        this.state.subscription_thank_you_message = this.props.user.subscriptions.thank_you
        this.state.headlineImageID = this.props.user.subscriptions.headline

        this.saveChanges = this.saveChanges.bind(this)


        this.setHeadlineImage = (file) => {
            this.setState({headlineImageID: file._id})
            console.log(file)
        }
    }

    saveChanges(){
        v1.user.updatePageConfigurables({
            supporter: this.state.supporter,
            subheadline: this.state.headlineImageID,
            thanks: this.state.thank_you_message,
            item: this.state.donation_item,
            price: parseFloat(this.state.donation_item_unit_price),
            supporters: this.state.allow_supporters,
            subthanks: this.state.subscription_thank_you_message,
            subsactive: this.state.subscriptions_active ,
            subprice: parseFloat(this.state.subscription_price),
            subperiod: parseInt(this.state.subscription_period),
            gacode: this.state.google_analytics_code,
        }).then(resp => {
            if(resp.status === 'saved'){
                alert('Successfully updated your preferences')
            }else {
                alert(resp.error)
            }
        }).catch(alert)
    }

    render(){
        return (
            <div class="box">
                <div class="box-body">
                    <div className="row justify-content-center">
                        <div class="col-md-11">

                            
                            {this.props.supporters ?
                                <>
                                    <div className="form-group">
                                        <div className="c-inputs-stacked">
                                            <input type="checkbox" id="checkbox_348" checked={this.state.allow_supporters} onChange={(evt) => { console.log(evt.target.value); this.setState({ allow_supporters: evt.target.checked }) }} />
                                            <label for="checkbox_348" className="block" onChange={(evt) => { this.setState({ allow_supporters: evt.target.checked }) }}>Allow {this.state.supporter}s to buy my a {this.state.donation_item}</label>
                                        </div>
                                    </div>
                                    <div class="form-group">
                                        <label>Call my supporters</label>
                                        <select class="form-control" onChange={(evt) => { this.setState({ supporter: evt.target.value }) }}  >
                                            {
                                                this.props.validSupporterNames.map(name => {
                                                    return <option value={name} selected={(name === this.state.supporter)} >{name}s</option>
                                                })
                                            }
                                        </select>
                                    </div>
                                    
                                    <div class="form-group">
                                        <label>My {this.state.supporter}s can buy me a </label>
                                        <select class="form-control" onChange={(evt) => { this.setState({ donation_item: evt.target.value }) }}  >
                                            {
                                                this.props.validDonationItems.map(name => {
                                                    return <option value={name} selected={(name === this.state.donation_item)} >{name}</option>
                                                })
                                            }
                                        </select>
                                    </div>
                                    <div class="form-group">
                                        <label class="control-label">Price of one {this.state.donation_item}</label>
                                        <PriceButtonGroup prices={[0.5, 1, 2, 5, 10]} price={this.state.donation_item_unit_price} onChange={price => this.setState({donation_item_unit_price: price})} />
                                    </div>
                                    <div class="form-group">
                                        <label>Thank you message</label>
                                        <p><small>This message is shown to your subscribers after they make a purchase or buy you a {this.state.donation_item}.</small></p>
                                        <RichEditor text={this.state.thank_you_message} onChange={(text) => { this.setState({ thank_you_message: text }) }} />
                                    </div>
                                </> : <></>}

                            
                            { this.props.subscriptions ? 
                            <>
                            <div className="form-group">
                                <label>Subscriptions</label>
                                <div className="c-inputs-stacked">
                                    <input type="checkbox" id="checkbox_347" checked={this.state.subscriptions_active} onChange={(evt) => { this.setState({ subscriptions_active: evt.target.checked }) }} />
                                    <label for="checkbox_347" className="block" onChange={(evt) => { this.setState({ subscriptions_active: evt.target.checked }) }} >Enable {this.state.supporter}s to subscribe to your content</label>
                                </div>
                            </div>
                            <div class="row">
                                <div class="col-md-6">
                                    <div className="form-group">
                                        <label>Subscription headline image</label>
                                        <div class="">
                                            <InlineUpload type="image" purpose="headline" onUploaded={this.setHeadlineImage} allowedTypes={["image/*"]} />
                                            <button class="btn btn-default btn-block " onClick={this.handleOpenF} ><i class="fa fa-upload"></i> Upload  image</button>
                                        </div>
                                    </div>
                                </div>
                                <div class="col-md-6">
                                    <div class="form-group">
                                        <label>Subscriptions last for: </label>
                                        <select class="form-control" value={this.state.subscription_period} onChange={(evt) => { this.setState({ subscription_period: evt.target.value }) }}  >
                                            <option value={3}>3 months</option>
                                        </select>
                                    </div>
                                    <div class="form-group">
                                                <label class="control-label">Price for 3 months subscription.</label>
                                                <PriceButtonGroup prices={[2,5,10,25,50]} price={this.state.subscription_price} onChange={price => this.setState({subscription_price: price})} />
                                    </div>
                                    <div class="form-group">
                                        <label>New Subscriber Thank you message</label>
                                        <p><small>This message is shown to your new subscribers. You can include private info here. i.e link to a private group</small></p>
                                        <RichEditor text={this.state.subscription_thank_you_message} onChange={(text) => { this.setState({ subscription_thank_you_message: text }) }} />
                                    </div>
                                </div>
                            </div>
                            
                            </> : <></>}
                            
                        </div>
                        
                    </div>
                    <div class="row justify-content-center" >
                        <div class="col-md-4" >
                            <button class="btn btn-block btn-primary" onClick={this.saveChanges} ><i class="fa fa-check"></i> Save changes</button>
                        </div>
                    </div>
                </div>
            </div>
        )
    }
}

class SettingsLandingPage extends Component {
    
    listItem(page) {
        return (
            <Link to={"/creator/settings/" + page.id} onClick={this.props.reloadFn} class="media-list bb-1 bb-dashed border-light">
                <div class="media align-items-center">
                    <div class={"avatar avatar-lg " + page.lit ? "status-success" : ""} >
                        <i class={"fa fa-"+ page.icon}></i>
                    </div>
                    <div class="media-body">
                        <p class="font-size-16">
                            <div class="hover-primary" ><strong>{page.title}</strong></div>
                        </p>
                    </div>
                </div>
                <div class="media pt-0 d-none d-md-block">
                    <p>{page.description}.</p>
                </div>
            </Link>)
    }

    render() {
        return(
            <div class="box">
                <div class="box-header with-border">
                    <h4 class="box-title">Settings page</h4>
                </div>
                <div class="box-body p-0">
                    {this.props.pagesFn().map(page => {
                        return page.verifiedOnly ? (this.props.user.verified ? this.listItem(page) : <></>) : this.listItem(page) 
                    })}
                </div>
            </div>
        )
    }
}



class UserSettings extends Component {
    constructor(props) {
        super(props)

        this.state = {
            loading: false,
            user: this.props.user,
        }

        this.reloadSettings = () => v1.user.current(false).then(user => {
            this.setState({loading: false, user})
        })

        this.pages = () => [
            {title: "My page", icon: "edit", lit: false, verifiedOnly: true, id: "mypage", 'description': 'Your page is what you share to the world. Edit how its displayed, your username and more.'},
            { title: "My subscriptions", icon: "heart text-danger", lit: false, verifiedOnly: true, id: "subscriptions", 'description': 'Allow paying user to subscribe o your private content for a specified time period.' },
            {title: "Buy me a "+this.state.user.page.donation_item, icon: "coffee text-success", lit: false, verifiedOnly: true, id: "buymeacoffee", 'description': 'Allow your supporters to support you'},
            { title: "Bank details and cashout", icon: "credit-card", lit: false, verifiedOnly: true, id: "payment", 'description': 'Add your bank details to cashout to your bank account.' },
            { title: "Personal settings", icon: "user text-info", lit: false, verifiedOnly: false, id: "personal", 'description': 'Change your personal settings.' },
            {title: "Password & security", icon: "lock text-warning", lit: false, verifiedOnly: false, id: "security", 'description': 'Change your password and security settings.'},
            { title: "Notifications", icon: "envelope", lit: false, verifiedOnly: false, id: "notifications", 'description': 'Fine tune your notifications.' },
            {title: "My Links", icon: "link", lit: false, verifiedOnly: true, id: "mylinks", 'description': 'Links to your page and more.'},
        ]

    }
    render() {
        return this.state.loading ? <Preloader /> : (
            <>
                
                <div class="mr-auto">
                    <div class="d-inline-block align-items-center">
                        <nav>
                            <ul class="breadcrumb fa-2x">
                                <li class="breadcrumb-item"><Link to="/creator/settings"><i class="fa fa-cog"></i> Settings</Link>
                                    <Switch>
                                        <Route path="/creator/settings/mypage">/ Page</Route>
                                        <Route path="/creator/settings/buymeacoffee">/ Buy me a {this.state.user.page.donation_item} </Route>
                                        <Route path="/creator/settings/subscriptions">/ Subscriptions</Route>
                                        <Route path="/creator/settings/personal" >/ Personal settings</Route>
                                        <Route path="/creator/settings/security">/ Security</Route>
                                        <Route path="/creator/settings/notifications">/ Notifications</Route>
                                        <Route path="/creator/settings/payment">/ Payments</Route>
                                        <Route path="/creator/settings/mylinks">/ Links</Route>
                                        
                                    </Switch>
                                </li>
                            </ul>
                        </nav>
                    </div>
                </div>
                <div class="tab-content tabcontent-border">
                    <div class="tab-pane active" id="activeTab" role="tabpanel">
                        <div class="box-body">
                            <Switch >
                                <Route path="/creator/settings/mypage" render={(props) => <PageSettings user={this.state.user} />} />
                                <Route path="/creator/settings/buymeacoffee" render={(props) => <AdvancedPageSettings
                                    user={this.state.user}
                                    supporters={true}
                                    validSupporterNames={["supporter", "fan", "client", "member"]}
                                    validDonationItems={["coffee", "beer", "lunch", "pizza", "chocolate", "candy", "ice cream", "airtime"]} />} />
                                <Route path="/creator/settings/subscriptions" render={(props) => <AdvancedPageSettings
                                    user={this.state.user}
                                    subscriptions={true}
                                    validSupporterNames={["supporter", "fan", "client", "member"]}
                                    validDonationItems={["coffee", "beer", "lunch", "pizza", "chocolate", "candy", "ice cream", "airtime"]} />} />
                                <Route path="/creator/settings/personal" render={(props) => <PersonalDetailsSettings user={this.state.user} />} />
                                <Route path="/creator/settings/security" render={(props) => <SecuritySettings user={this.state.user} />} />
                                <Route path="/creator/settings/notifications" render={(props) => <NotificationSettings user={this.state.user} />} />
                                <Route path="/creator/settings/mylinks" render={(props) => <MyLinks user={this.state.user} />} />
                                <Route path="/creator/settings/payment" render={(props) => <Wallet hideRecent={true} user={this.state.user} hideWithdraw={true} />} />
                                <Route path="/creator/settings" render={(props) => <SettingsLandingPage user={this.state.user} reloadFn={this.reloadSettings} pagesFn={this.pages} />} />
                            </Switch>
                        </div>
                    </div>
                </div>
        </>
        )
    }
}

export default UserSettings;