
import React, { Component } from "react"
import BottomNav from "./BottomNavBar";
import SubscriptionPayModal from "./payments/SubscriptionPayModal";
import SupportCreatorModal from "./payments/SupportCreatorModal";
import SharePageModal from "./SharePageModal";


class CreatorBottomNav extends Component {
    constructor(props) {
        super(props)

        this.closeMenu = () => {
            document.getElementById('mobileNavBtn').click()
        }

        this.state = {
            showBuyModal: false,
            showShareModal: false,
            showSubscribeModal: false,
        }

        this.sharePage = this.sharePage.bind(this)
        this.buyMeACoffee = this.buyMeACoffee.bind(this)
        this.subscribe = this.subscribe.bind(this)
    }

    async sharePage(evt) {
        try {
            await navigator.share({
                title: 'Check out this post ',
                text: "I think you'll find it interesting. Check it out!",
                url: window.location.href,
            }).then(_ => {
                console.log('ok')
                this.setState({ showShareModal: false })
            }).catch(err => {
                console.log(err)
                this.setState({ showShareModal: true })
            })    
        } catch (e) {
            this.setState({ showShareModal: true })
        }
    }

    buyMeACoffee(evt) {
        this.setState({ showBuyModal: !this.state.showBuyModal })
        //return (evt && evt.preventDefault ? evt.preventDefault() : null)
    }

    subscribe(evt) {
        this.setState({ showSubscribeModal: !this.state.showSubscribeModal })
        //return (evt.preventDefault ? evt.preventDefault() : null)
    }

    render() {
        return this.props.user._id === this.props.creator._id ? <BottomNav user={this.props.creator} /> : (
<nav class="navbar fixed-bottom navbar-expand-md custom-navbar navbar-container d-xs-block d-sm-block d-md-none" style={{minHeight: 'unset'}}>
    <div class="container">
        {this.state.showBuyModal ? <SupportCreatorModal user={this.props.user} content={this.props.content} creator={this.props.creator} amount={this.props.creator.page.donation_item_unit_price} items={1} itemName={this.props.creator.page.donation_item} onClose={this.buyMeACoffee} /> : <></> }
        {this.state.showShareModal ? <SharePageModal onClose={_ => this.setState({showShareModal: false})} /> : <></>}
        {this.state.showSubscribeModal ? <SubscriptionPayModal user={this.props.user} content={this.props.content} creator={this.props.creator} amount={this.props.creator.subscriptions.price}  onClose={_ => this.setState({showSubscribeModal: false})} /> : <></>}
        
        <div class="mobile-app-icon-bar bg-dark" style={{'display': 'flex', 'width': '100%'}}>
                        {this.props.creator.page.allow_supporters ?
                            <a href="javascript:;" className="btn btn-info btn-block" onClick={this.buyMeACoffee} ><i class="fa fa-coffee" aria-hidden="true"></i> Buy me a {this.props.creator.page.donation_item } </a> : <></>}
                        {(this.props.creator.subscriptions.active &&  !this.props.creator.page.allow_supporters) ?
                            <a href="javascript:;" className="btn btn-info btn-block" onClick={this.subscribe} ><i class="fa fa-lock" aria-hidden="true"></i> Subscribe to view exclusive content</a>
                        : <></>}
            <a href="javascript:;"  className="btn" onClick={this.sharePage} ><i class="fa  fa-share-alt" aria-hidden="true"></i></a>
            
        </div>
    </div>
</nav>
        );
    }
}

/*
            
*/

export default CreatorBottomNav;

