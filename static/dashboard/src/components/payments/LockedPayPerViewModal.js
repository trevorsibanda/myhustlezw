
import React, {Component} from 'react'
import SweetAlert from 'react-bootstrap-sweetalert';
import PaymentUI from './PaymentUI';

class LockedPayPerViewModal extends Component{ 
    
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
            <SweetAlert show={this.state.show} showConfirm={false} showCancel={false} showCloseButton={true} title={'Pay to unlock content.'} onCancel={this.onClose}  >
                <PaymentUI purpose='pay_per_view' onClose={this.onClose} content={this.props.content} target={this.props.target} creator={this.props.creator} amount={this.props.content.price} user={this.props.user} />
            </SweetAlert>
        )
    }
}


export default LockedPayPerViewModal;

