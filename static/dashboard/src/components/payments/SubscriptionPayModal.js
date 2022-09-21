
import React, {Component} from 'react'
import SweetAlert from 'react-bootstrap-sweetalert';
import PaymentUI from './PaymentUI';

class SubscriptionPayModal extends Component{ 
    
    constructor(props) {
        super(props);
        this.state = {
            show: true,
        }

        this.onClose = this.onClose.bind(this)
    }

    onClose() {
        this.setState({
            show: false
        })
        return this.props.onClose ? this.props.onClose() : null
    }

    render() {
        return (
           <SweetAlert show={this.state.show} showConfirm={false} showCancel={false} showCloseButton={true} title={'Subscribe to @'+this.props.creator.username} onCancel={this.onClose}  >
              <PaymentUI purpose="subscribe" onClose={this.onClose} content={this.props.content} target={this.props.target} creator={this.props.creator} amount={this.props.creator.subscriptions.price} user={this.props.user} />
            </SweetAlert>
        )
    }
}


export default SubscriptionPayModal;

