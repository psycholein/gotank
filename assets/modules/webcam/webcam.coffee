window.Webcam = class Webcam extends window.BasicModule
  afterInit: =>
    img = "<image src='http://#{document.location.hostname}:8080/?action=stream' />"
    $module = $(@selector).filter("[data-name='#{@name}']")
    $module.find('.img').html(img)

  config: =>
    @position = "left"
