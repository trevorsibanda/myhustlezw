import React, { Component } from "react";
import { Link } from "react-router-dom"
import v1 from "../api/v1";



class EnableSubscriptionsPrompt extends Component {
    constructor(props){
        super(props)
    }
    render() {
        return (
            <>
                <div className="box justify-content-center">
                    <img src="/assets/img/bg/social/bg.png" className="img" />
                    <h4 className="title padding-top-40">Hey {this.props.user.fullname}, enable memberships to increase earnings.</h4>
                    <p>Subscription groups allow you to create content exclusively for your {this.props.user.profile.supporters}.
                            <br />Your fans pay a subscription to gain access to one or more campaigns.
                            <br />When you create a new campaign you can choose its subscription group and only paid members of that group will be able to view the content.
                            <strong>Currently you are limited to one subscription group.</strong></p>
                </div>
                <div className="row justify-content-center">
                    <div className="col-md-4">
                        <button onClick={this.props.enableSubscriptions} className="btn btn-danger btn-block"><i className="fa fa-heart"></i> Enable subscription accounts</button>
                    </div>
                </div>
            </>
        )
    }
}

class SubscriptionsManager extends Component {
    render(){
        return (
            <div class="row">
                <h1>Enable subscriptions</h1>
            </div>
        )
    }
}

class Subscriptions extends Component {

    constructor(props) {
        super(props)
        this.state = {
            user: {
                fullname: 'Creator!',
                memberships: false,
                payments_active: true,
                profile: {
                    supporters: 'supporter',
                }
            }
        }
        v1.user.current().then(user => {
            this.setState({ user, ready: true })
        }).catch(err => {
            //alert of err
        })

        this.enableSubscriptions = this.enableSubscriptions.bind(this)
    }

    enableSubscriptions() {
        v1.subscriptions.enable().then(enable => {
            if(enable.status === true){
                this.setState({user: enable.user})
            }else {
                alert(enable.error)
            }
        })
    }

    render() {
        return (
            <>
                
                <div class="box">
                    <div class="box-header">
                        <h4 class="box-title">Memberships</h4>
                    </div>
                    <div class="box-body">
                        {this.state.user.memberships ? <SubscriptionsManager user={this.state.user} /> : <EnableSubscriptionsPrompt enableSubscriptions={this.enableSubscriptions} user={this.state.user} /> }
                    </div>
                </div>
                
                
            </>
        )
    }

}

export default Subscriptions;