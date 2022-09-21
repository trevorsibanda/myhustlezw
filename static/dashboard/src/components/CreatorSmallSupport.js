import React, { Component } from "react"
import { Link, Switch, Route } from 'react-router-dom'
import OptionsButtonGroup from "./OptionsButtongroup"
import SupportCreatorModal from "./payments/SupportCreatorModal"
import money from "./payments/Amount"

class CreatorSmallSupport extends Component{

    constructor(props) {
        super(props)
        this.state = {
            buyItems: 1,
            openPayModal: false,
        }
        this.itemOptions = [1, 2, 3, 5, 10]
        this.setBuyItems = this.setBuyItems.bind(this)
        this.payModal = this.payModal.bind(this)
    }

    setBuyItems(n) {
        this.setState({
            buyItems: n
        })
    }

    payModal() {
        this.setState({openPayModal: !this.state.openPayModal})
    }


    render() {
        let styleActive = "activesetter"
        return this.state.openPayModal ?
            <SupportCreatorModal onClose={this.payModal} user={this.props.user} creator={this.props.creator} items={this.state.buyItems} itemName={this.props.creator.page.donation_item} amount={this.props.creator.page.donation_item_unit_price * this.state.buyItems} /> : (
            <div class="box">
                    <div class="box-body">
                        <div class="flexbox align-items-baseline mb-20">
                            <h6 class="text-uppercase ls-2">Buy {this.props.creator.fullname} a {this.props.creator.page.donation_item}</h6>
                            <small>{money.formatUSD(this.props.creator.page.donation_item_unit_price)} each</small>
                        </div>
                        <div class=" gap-y padding-top-10">
                            <strong>Buy {this.props.creator.page.donation_item} x</strong>
                            <OptionsButtonGroup item={this.state.buyItems} items={this.itemOptions} onChange={this.setBuyItems} />
                        <button class="btn mt-5 btn-info btn-block pull-right" onClick={this.payModal}>Pay {money.formatUSD(this.state.buyItems *  this.props.creator.page.donation_item_unit_price)}</button>
                        </div>
                    </div>
                    <div class="box-footer">You can pay using Ecocash for Zimbabwe payments and through 2Checkout for International USD
                        payments.</div>
                </div>
        )
    }
}


export default CreatorSmallSupport;