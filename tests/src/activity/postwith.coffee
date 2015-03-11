helpers = require '../helpers/helpers.js'
assert  = require 'assert'

activitySelector = '[testpath=activity-list] section:nth-of-type(1) [testpath=ActivityListItemView]:first-child'

module.exports =


  postMessageWithCode: (browser) ->

    helpers.beginTest(browser)

    timestamp = Date.now()
    code      = "console.log('#{timestamp}')"
    post      = '```' + code + '```'
    selector  = '[testpath=ActivityListItemView]:first-child .has-markdown code'

    helpers.doPostActivity(browser, post, no)

    browser
      .assert.containsText selector, code # Assertion
      .end()


  postMessageWithImage: (browser) ->

    helpers.beginTest(browser)

    image = 'http://placehold.it/200x100'
    selector = activitySelector + ' .activity-content-wrapper .embed-image-view img'

    browser
      .click                  '[testpath="public-feed-link/Activity/Topic/public"]'
      .waitForElementVisible  '[testpath=ActivityInputView]', 25000
      .click                  '[testpath="ActivityTabHandle-/Activity/Public/Recent"]'
      .click                  '[testpath=ActivityInputView]'
      .setValue               '[testpath=ActivityInputView]', image
      .pause                  5000
      .click                  '[testpath=post-activity-button]'
      .waitForElementVisible   selector, 20000 # Assertion
      .end()


  postMessageWithLink: (browser) ->

    helpers.beginTest(browser)

    link = 'http://google.com/' # last '/' is the trick!
    linkSelector = activitySelector + ' .activity-content-wrapper article a'

    helpers.doPostActivity(browser, link, yes, yes)

    browser.getAttribute linkSelector, 'href', (result) ->
      href = result.value
      assert.equal(link, href)

    browser.end()


  postCommentWithCode: (browser) ->

    helpers.beginTest(browser)

    post      = helpers.getFakeText()
    timestamp = Date.now()
    code      = "console.log(#{timestamp})"
    comment   = '```' + code + '```'
    selector  = '[testpath=ActivityListItemView]:first-child .comment-contents .has-markdown code'

    helpers.doPostActivity(browser, post)
    helpers.doPostComment(browser, comment, no)

    browser
      .pause 2000
      .assert.containsText selector, code # Assertion
      .end()


  postCommentWithImage: (browser) ->

    helpers.beginTest(browser)

    post     = helpers.getFakeText()
    image    = 'http://placehold.it/200x100'
    selector = activitySelector + ' .comment-contents .comment-body-container .has-markdown a'

    helpers.doPostActivity(browser, post)
    helpers.doPostComment(browser, image, yes, yes)

    browser
      .assert.containsText selector, image # Assertion
      .end()


  postCommentWithLink: (browser) ->

    helpers.beginTest(browser)

    post     = helpers.getFakeText()
    link = 'http://google.com/' # last '/' is the trick!
    linkSelector = activitySelector + ' .comment-contents .comment-body-container .has-markdown a'

    helpers.doPostActivity(browser, post)
    helpers.doPostComment(browser, link, yes, yes)

    browser.getAttribute linkSelector, 'href', (result) ->
      href = result.value
      assert.equal(link, href)

    browser.end()
