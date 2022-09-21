import { Component } from "react";
import { Link } from "react-router-dom"
import ReactTimeAgo from "react-time-ago/commonjs/ReactTimeAgo";
import v1 from "../api/v1";
import Preloader from "../components/PreLoader";
import RichEditor from "../components/RichEditor"
import money from "../components/payments/Amount"

class ViewSupporter extends Component{

  constructor(props) {
    super(props)

    v1.page.set({ title: 'View supporter' })
    v1.page.track()

    this.state = {
      support: this.props.support ? this.props.support : {},
      loading: true,
      refund_reason: '',
    }

    let {id} = this.props.match.params

    if (this.state.support._id) {
      this.setState({loading: false})
      
    } else {
      v1.supporters.get(id, false).then(support => {
        this.setState({ support, loading: false })
      }).catch(err => {
        this.setState({ loading: false })
      })
    }

    this.hidePublicMessage = this.hidePublicMessage.bind(this)
    this.sendPrivateEmail = this.sendPrivateEmail.bind(this)
    this.hideActivity = this.hideActivity.bind(this)
    this.fulfillOrder = this.fulfillOrder.bind(this)
    this.refundOrder = this.refundOrder.bind(this)

  }

  fulfillOrder() {
    this.setState({loading: true})
    v1.supporters.modifyServiceOrderStatus('fulfill', this.state.support._id, this.state.support).then(resp => {
      if (resp.error) {
        alert('Failed with error\n' + resp.error)
        return
        }
      this.setState({ loading: false, support: resp })
      alert('Great job :)\nYou have completed this order.')
    }).catch(err => {
      alert(err && err.error ? err.error : 'Failed to fulfill order')
      this.setState({loading: false})
    })
  }

  refundOrder() {
    if (this.state.refund_reason.length < 10) {
      alert('Reason for refunding order must be at least 10 characters long')
      return
    }
    this.setState({loading: true})
    v1.supporters.modifyServiceOrderStatus('refund', this.state.support._id, {reason: this.state.refund_reason}).then(resp => {
      if (resp.error) {
        alert('Failed with error\n' + resp.error)
        return
        }
        this.setState({loading: false, support: resp})
        alert('You have cancelled and refunded ths order. Details sent to payer via email')
    }).catch(err => {
      alert(err && err.error ? err.error : 'Failed to refund order')
      this.setState({loading: false})
    })
  }

  hidePublicMessage = () => {
    v1.supporters.hideActivity(this.state.support._id, "message").then(support => {
      this.setState({ support })
      alert(support.hide_message ? "Message now hidden" : "Message now showing")
    })
  }

  hideActivity = () => {
    v1.supporters.hideActivity(this.state.support._id, "all").then(support => {
      this.setState({ support })
      alert(support.hidden ? "This will no longer appear on your public page" : "Activity now showing on public page")
    })
  }

  sendPrivateEmail = () => {

    alert("Sorry this feature is not yet available.")
  }


  render() {
      var message

    if (!this.state.loading) {
      switch (this.state.support.support_type) {
        case "subscribed": message = " subscribed to @" + this.props.creator.username + "'s account";  break;
          case "paid_content": message =  " paid to view " + this.state.support.item_name; break;
          case "support": message = " bought you " + this.state.support.items + " " + this.state.support.item_name; break;
          case "service_request": message = " placed an order for '"+ this.state.support.item_name+"'"; break;
          default: message = "contributed somehow";
      }
    }
    
    return  this.state.loading ? <Preloader /> : (this.state.support._id ? 
<>
<h4 class="margin-top-20">{ this.state.support.display_name }</h4>
<hr class="hr-danger" />
<div class="card" >
  <img class="card-img-top" src="/assets/img/placeholder.png" alt="Card cap" />
  <div class="card-body">
      <h5 class="card-title">{message}</h5>
      <p class="card-text">For <b>{ money.format(this.state.support.currency, this.state.support.amount) }</b></p>
  </div>
  <ul class="list-group list-group-flush">
    <li class="list-group-item"><ReactTimeAgo date={this.state.support.created_at}/></li>
    <li class="list-group-item">{this.state.support.created_at} </li>
    <li class="list-group-item">ID: {this.state.support._id}</li>
      {this.state.support.form ?
        <li class="list-group-item">
          <table class="table table-responsive table-striped">
            <tbody>
              <tr>
                <td>Order for </td>
                <td><Link to={"/creator/content/"+this.state.support._id}>{this.state.support.item_name}</Link></td>

              </tr>
              <tr>
                <td>Fullname</td>
                <td>{this.state.support.form.fullname}</td>
              </tr>
              <tr>
                <td>Email</td>
                <td>{this.state.support.form.email}</td>
              </tr>
              <tr>
                <td>Phone number</td>
                <td>{this.state.support.form.phone}</td>
              </tr>
              {this.state.support.form.fulfilled ?
              <tr>
                <td>Fulfilled </td>
                <td><ReactTimeAgo date={this.state.support.form.fulfilled_at}/> {this.state.support.form.fulfilled_at}</td>
              </tr> : <></>}
              {this.state.support.form.refunded ?
              <tr>
                <td>Refunded  </td>
                <td><ReactTimeAgo date={this.state.support.form.refunded_at}/> {this.state.support.form.refunded_at}</td>
              </tr> : <></>}
              <tr>
                <td>Instructions</td>
                <td>{this.state.support.form.instructions}</td>
              </tr>
              <tr>
                <td>{this.state.support.form.question}</td>
                <td>{this.state.support.form.answer}</td>
              </tr>
              <tr>
                <td>Items left</td>
                <td>{this.state.support.form.quantity_left}  at time of ordering</td>
              </tr>
            </tbody>
          </table>
        </li> :
        <li class="list-group-item"><b>Message from supporter</b>
          <textarea class="form-control" rows={10} readOnly={true} placeholder="Message from suppporter" value={this.state.support.comment}>
          </textarea>
        </li>}

  </ul>
          <div class="card-body">
            {this.state.support.form && !this.state.support.form.refunded ?
              <p>
                Mark this order as completed. <b>Only mark it as complete after you've fulfilled the required task.</b> We will immediately send email notifications to both you and the person who made the order.
                <button onClick={this.fulfillOrder} disabled={this.state.support.form.fulfilled} class="card-link btn btn-block btn-warning mb-20">
                  <i class="fa fa-check" ></i> {this.state.support.form.fulfilled ?  "You have fulfilled this task" : "Mark this task as fulfilled"}
                </button>
              </p> : <></>}
            <p>
              You can hide this activity from appearing on your public profile page.
              <button onClick={this.hideActivity}  class="card-link btn btn-block btn-warning mb-20">
                <i class="fa fa-eye" ></i> {this.state.support.hidden ? "Show this on your public page" : "Hide activity from page"}
              </button>
            </p>
            <p>
              You can hide the message left from been shown on your page.
              <button onClick={this.hidePublicMessage}  class="card-link btn btn-block btn-info mb-20">
                <i class="fa fa-eye" ></i> {this.state.support.hide_message ? "Show message on your public page" : "Hide message on your public page"}
              </button>
            </p>
            <p>
              You can send a private email message to this person.
              <button onClick={this.sendPrivateEmail} class="card-link btn btn-block btn-success mb-20">
                <i class="fa fa-envelope" ></i>  Send private email
              </button>
            </p>
            {this.state.support.form && !this.state.support.form.fulfilled ?
            <p>
                Cancel and refund this order. We will handle the refund to the person's account, <b>many refund requests from your account will get your account deactivated.</b>
              <div className="form-group">
                <label>Reason for cancelling</label>
                <div class="input-group mb-3">
                    <textarea maxLength={1024} rows={5} onChange={(evt) =>this.setState({ refund_reason: evt.target.value }) }  class="form-control" placeholder="reason for cancelling" >
                    {this.state.refund_reason}
                    </textarea> 
                </div>
            </div>
              <button onClick={this.refundOrder} disabled={this.state.support.form.refunded} class="card-link btn btn-block btn-danger mb-20">
                <i class="fa fa-warning" ></i> {this.state.support.form.refunded ? "Order cancelled and refunded" : "Cancel and refund this order"}
              </button>
            </p> : <></>}
    
  </div>
</div>
</> : <div class="alert alert-danger">Failed to load this supporter. An error occured.</div>
        )
    }
}

export default ViewSupporter;