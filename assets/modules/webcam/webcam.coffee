window.Webcam = class Webcam extends window.BasicModule
  afterInit: =>
    setTimeout (-> @getPicture), 1000

  getPicture: =>
    console.log('get picture')
    @event.send(@module, @name, "get", {})

  config: =>
    @position = "left"

  router: (event) =>
    switch event.Task
      when "img"
        img = "<image src='#{event.Data.img}' />"
        $module = $(@selector).filter("[data-name='#{event.Name}']")
        $module.find('.img').html(img)
