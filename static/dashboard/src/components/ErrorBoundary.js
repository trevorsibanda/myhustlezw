import { Component } from "react";
import { Redirect } from "react-router-dom";
import { Link } from "react-router-dom";


class ErrorBoundary extends Component {

    constructor(props) {
        super(props)
        this.state = {
            hasError: false,
            error: null,
        }
    }

    componentDidCatch(error, _) {
        this.setState({
            hasError: true,
            error: error,
        })
        alert({
            title: 'Error loading this page',
            text: 'An error occured trying to load this page.\n\n' + (error.stack ? error.stack : JSON.stringify(error)),
            icon: 'error',
            timer: 10000,
        })

        setTimeout(_ => {
            this.setState({ hasError: false, error: null })
        }, 250)

    }

    componentDidUpdate(prevProps) {
        console.log(this.props, prevProps)
        if (this.props.children !== prevProps.children) {
            this.setState({ hasError: false, error: null })
        }
    }

    render() {
        return this.state.hasError ? <Redirect to={this.props.recoverPage ? this.props.recoverPage : "/creator/dashboard"} /> : this.props.children
    }
}

export default ErrorBoundary;