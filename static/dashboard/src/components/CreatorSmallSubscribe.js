import React, { Component } from "react"
import money from "./payments/Amount"
import SubscriptionPayModal from "./payments/SubscriptionPayModal"

class CreatorSmallSubscribe extends Component{

    constructor(props) {
        super(props)
        this.state = {
            buyItems: 1,
            openPayModal: false,
        }
        this.payModal = this.payModal.bind(this)
    }

    payModal() {
        this.setState({openPayModal: !this.state.openPayModal})
    }


    render() {
        return this.state.openPayModal ?
            <SubscriptionPayModal onClose={this.payModal} user={this.props.user} creator={this.props.creator} items={1} amount={this.props.creator.subscriptions.price} target={window.location.href} /> : (
            <div class="box">
                    <div class="box-body">
                        <div class="flexbox align-items-baseline mb-20">
                            <p>
                                Gain access to exclusive content and subscribe right now. <br />
                                Your subscription is instant.
                            </p>
                        </div>
                        <div class=" gap-y padding-top-10">
                        
                        <button class="btn mt-5 btn-info btn-block pull-right" onClick={this.payModal}>Subscribe now for {money.formatUSD(this.props.creator.subscriptions.price)}</button>
                        </div>
                    </div>
                    <div class="box-footer">You can pay using Ecocash for Zimbabwe payments and through 2Checkout for International USD
                        payments.</div>
                </div>
        )
    }
}


export default CreatorSmallSubscribe;