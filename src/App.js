import React, { Component } from 'react';
import { Grid, Navbar, Button, Form, FormGroup, InputGroup, FormControl, ListGroup, ListGroupItem, Glyphicon} from 'react-bootstrap';

class App extends Component {
  constructor(props) {
    super(props)
    this.state = {
      value: "",
      results: []
    }
    this.handleChange = this.handleChange.bind(this)
    this.handleSearch = this.handleSearch.bind(this)
  }

  handleChange(e) {
    this.setState({ value: e.target.value })
  }

  handleSearch(e) {
    e.preventDefault()
    let query = this.state.value.trim()
    if (query.length === 0) {
      return
    }

    fetch(`https://pt-search-backend-4ywt76ki5q-uc.a.run.app/api/search?q=${query}`)
      .then((response) => response.json())
      .then((responseJson) => {
        this.setState({
          results: responseJson
        })
      })
      .catch((error) => {

      })
  }

  render() {
    return (
      <div>
        <Navbar>
          <Navbar.Header>
            <Navbar.Brand>
              <a href="/">PT-Search</a>
            </Navbar.Brand>
          </Navbar.Header>
        </Navbar>
        <Grid>
          <Form onSubmit={this.handleSearch}>
            <FormGroup>
              <InputGroup>
                <FormControl
                  type="text"
                  value={this.state.value}
                  placeholder="Type to search"
                  onChange={this.handleChange}
                />
                <InputGroup.Button>
                  <Button type="submit" onClick={this.handleSearch}>Search</Button>
                </InputGroup.Button>
              </InputGroup>
            </FormGroup>
          </Form>
          <ResultList results={this.state.results}></ResultList>
        </Grid>
      </div>
    );
  }
}

class MovieItem extends Component {
  constructor(props) {
    super(props);
    this.state = {
      downloadStatus: ""
    }
    this.handleDownload = this.handleDownload.bind(this)
  }

  handleDownload(e) {
    let apiUrl = `https://pt-search-backend-4ywt76ki5q-uc.a.run.app/api/queue?from=${this.props.from}&id=${this.props.id}`
    console.log('send download request: ' + apiUrl)

    this.setState({downloadStatus: "downloading"})
    fetch(apiUrl, {method: 'POST'})
      .then((response) => {
        if (!response.ok) {
          this.setState({downloadStatus: "failed"})
          throw new Error('failed to send download request')
        } else {
          this.setState({downloadStatus: "finished"})
        }
      }).catch((err) => {
        console.log(err)
      })
  }

  render() {
    return <ListGroupItem>
      <div>
        <a href={this.props.url} target="_blank">
          {this.props.title}
        </a>
        <p>{this.props.sub_title}</p>
      </div>
      <div>
        <span>[{this.props.from}] </span>
        <span>{this.props.age} </span>
        <span>{this.props.size} </span>
        <span>{this.props.seeder}</span>
        <DownloadButton downloadStatus={this.state.downloadStatus} onClick={this.handleDownload}></DownloadButton>
      </div>
    </ListGroupItem>
  }
}

const DownloadButton = (props) => {
  switch(props.downloadStatus) {
    case 'downloading':
      return (
        <Button bsStyle="link" bsSize="sm" onClick={props.onClick} disabled>
          <Glyphicon glyph="sort-by-attributes" />
        </Button>
      )
    case 'finished':
      return (
        <Button bsStyle="link" bsSize="sm" disabled>
          <Glyphicon glyph="saved" />
        </Button>
      )
    default :
      return (
        <Button bsStyle="link" bsSize="sm" onClick={props.onClick}>
          <Glyphicon glyph="download-alt"/>
        </Button>
      )
  }
}

const ResultList = (props) => {
  const items = props.results
  const listItems = items.map((item) =>
    <MovieItem {...item} key={item.from + item.id}></MovieItem>
  )

  return <ListGroup>{listItems}</ListGroup>
}

export default App;
