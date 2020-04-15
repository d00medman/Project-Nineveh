from ctypes import *
import os

file_name = './nineveh.so'
prototype_file_name ='./nineveh_one_step.so'

def _load_nineveh():
    return cdll.LoadLibrary(file_name)


nineveh = _load_nineveh()

nineveh.NewNintendo.restype = c_longlong
nineveh.RunNintendo.argtypes = None
nineveh.RunNintendo.restype = None
nineveh.GetObservation.argtypes = None
nineveh.GetObservation.restype = c_char_p
nineveh.TakeAction.argtypes = [c_longlong]
nineveh.TakeAction.restype = c_float
nineveh.CloseEmulator.argtypes = None
nineveh.CloseEmulator.restype = None

nineveh.Reset.argtypes = None
nineveh.Reset.restype = None
nineveh.OpenToStart.argtypes = None
nineveh.OpenToStart.restype = None
nineveh.EndRecording.argtypes = None
nineveh.EndRecording.restype = None
nineveh.IsGameOver.argtypes = None
nineveh.IsGameOver.restype = c_int

class NinevehInterface(object):

    def __init__(self, use_gui=False, game_file='Castlevania.nes', frame_skip=16):
        # display modes
        mode = c_char_p(b'gui') if use_gui else c_char_p(b'headless')
        file = c_char_p(game_file.encode('utf-8'))

        nineveh.NewNintendo(file, mode, frame_skip)

    def start(self):
        nineveh.RunNintendo()

    def act(self, action):
        reward = nineveh.TakeAction(action)
        print(f'Reward {reward} for action {action}')
        return reward

    def reset(self):
        nineveh.Reset()
        nineveh.OpenToStart()

    def save_run_recording(self):
        nineveh.EndRecording()

    def get_observation(self):
        o = nineveh.GetObservation()
        # print(f'Raw observation in python: {o}')
        # print(type(o.decode('utf-8')))
        return o.decode('utf-8')

    def is_game_over(self):
        return bool(nineveh.IsGameOver())

# Both will eventually be removed, primarily needed for testing purposes
def act_loop():
    n = NinevehInterface(use_gui=False)
    n.reset()
    for i in range(1000000):
        a = int(input("Select action: "))
        print(f'Action #{i}, selected: {a}')
        # print(f'Advance to frame #{i} (pid: {os.getpid()})')
        if a == 99:
            n.save_run_recording()
        else:
            n.act(a)
            n.get_observation()
            print(n.is_game_over())
            print('-----')

def standard_use():
    n = NinevehInterface()
    nineveh.OpenToStart()
    n.start()

def make_output_dir(dirName='output'):
    if not os.path.exists(dirName):
        os.mkdir(dirName)
        print("Directory ", dirName, " Created ")
    else:
        print("Directory ", dirName, " already exists")


if __name__ == '__main__':
    print(f'Running emulator, pid: {os.getpid()}')
    make_output_dir()
    # standard_use()
    act_loop()
