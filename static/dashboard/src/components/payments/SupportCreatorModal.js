
import React, {Component} from 'react'
import SweetAlert from 'react-bootstrap-sweetalert';
import PaymentUI from './PaymentUI';

class SupportCreatorModal extends Component{ 
    
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
            <SweetAlert show={this.state.show} showConfirm={false} showCancel={false} showCloseButton={true} title={'Buy @'+ this.props.creator.username + ' '+ this.props.items + ' '+ this.props.itemName + '(s)'} onCancel={this.onClose}  >
                <PaymentUI purpose='support' user={this.props.user} creator={this.props.creator} amount={this.props.amount}  items={this.props.items} itemName={this.props.itemName} content={this.props.content}  />
            </SweetAlert>
        )
    }
}


export default SupportCreatorModal;